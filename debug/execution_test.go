package UnitTests

import (
	"context"
	"flag"
	"fmt"
	"os"
	"strings"
	"testing"

	ce "github.com/fl4cko11/dbc-utility/internal/command_execution"
	cp "github.com/fl4cko11/dbc-utility/internal/command_processing"
	"github.com/fl4cko11/dbc-utility/logs"
	"github.com/pashagolub/pgxmock/v2"
)

func TestCommandExecution_Integration(t *testing.T) {
	testLogger := logs.InitLogger(os.Stderr)
	testLogger.ExitFunc = func(code int) {
		panic("fatal error occurred") // Чтобы после log.Fatal() функцией выхода была panic() а не ox.Exit (чтобы все тесты прогнать)
	}

	tests := []struct {
		name string
		args []string
	}{
		{
			name: "successful single DB removal",
			args: []string{"cmd", "-databases=testdb", "-operation=remove", "-pgpass=pass"},
		},
		{
			name: "successful multiple DBs removal",
			args: []string{"cmd", "-databases=db1,db2", "-operation=remove", "-pgpass=pass"},
		},
		{
			name: "template db delete",
			args: []string{"cmd", "-databases=db1%", "-operation=remove", "-pgpass=pass"},
		},
		{
			name: "single db and template",
			args: []string{"cmd", "-databases=db1,db2%", "-operation=remove", "-pgpass=pass"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil {
					t.Errorf("Ошибка: %v", r)
				}
			}()

			flag.CommandLine = flag.NewFlagSet("cmd", flag.ExitOnError) // Сбрасываем для кажого теста
			os.Args = tt.args                                           // Переопределяем аргументы коммандной строки для конкретного теста

			args := cp.CommandProcessing(testLogger)

			mock, err := pgxmock.NewConn(pgxmock.QueryMatcherOption(pgxmock.QueryMatcherEqual)) // обязательно включаю сравнение юнит запроса и мок запроса как строки
			if err != nil {
				t.Fatalf("Не удалось создать mock подключение: %v", err)
			}
			defer mock.Close(context.Background())

			for _, dbName := range args.DbNames {
				if strings.Contains(dbName, "%") {
					dbNameWoPersent := strings.Trim(dbName, "%")

					rows := mock.NewRows([]string{"datname"})
					rows.AddRow(fmt.Sprintf("%s_a", dbNameWoPersent))
					rows.AddRow(fmt.Sprintf("%s_b", dbNameWoPersent))

					mock.ExpectQuery(fmt.Sprintf("SELECT datname FROM pg_database WHERE datname LIKE '%s' AND datistemplate = false AND datname NOT IN ('postgres','template0','template1');", dbName)).WillReturnRows(rows)

					mock.ExpectExec(fmt.Sprintf("DROP DATABASE %s_a", dbNameWoPersent)).WillReturnResult(pgxmock.NewResult("DROP", 1))
					mock.ExpectExec(fmt.Sprintf("DROP DATABASE %s_b", dbNameWoPersent)).WillReturnResult(pgxmock.NewResult("DROP", 1))
				} else {
					mock.ExpectExec(fmt.Sprintf("DROP DATABASE %s", dbName)).WillReturnResult(pgxmock.NewResult("DROP", 1))
				}
			}

			ctx := context.Background()
			ce.CommandExecution(ctx, mock, args, testLogger)

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("Ожидаемые запросы не выполнены: %v", err)
			}
		})
	}
}
