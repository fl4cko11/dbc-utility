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

func TestCommandExecution(t *testing.T) {
	testLogger := logs.InitLogger(os.Stderr)
	testLogger.ExitFunc = func(code int) {
		panic("fatal error occurred") // Чтобы после log.Fatal() функцией выхода была panic() а не ox.Exit (чтобы все тесты прогнать)
	}

	tests := []struct {
		name        string
		args        []string
		expectFatal bool
	}{
		{
			name:        "successful single DB removal",
			args:        []string{"cmd", "-databases=testdb", "-operation=remove", "-pgpass=pass", "-debug"},
			expectFatal: false,
		},
		{
			name:        "successful single DB removal wo debug",
			args:        []string{"cmd", "-databases=testdb", "-operation=remove", "-pgpass=pass"},
			expectFatal: false,
		},
		{
			name:        "successful multiple DBs removal",
			args:        []string{"cmd", "-databases=db1,db2", "-operation=remove", "-pgpass=pass", "-debug"},
			expectFatal: false,
		},
		{
			name:        "template db removal",
			args:        []string{"cmd", "-databases=db1%", "-operation=remove", "-pgpass=pass", "-debug"},
			expectFatal: false,
		},
		{
			name:        "single db and template removal",
			args:        []string{"cmd", "-databases=db1,db2%", "-operation=remove", "-pgpass=pass", "-debug"},
			expectFatal: false,
		},
		{
			name:        "importand db [postgres] removal",
			args:        []string{"cmd", "-databases=postgres", "-operation=remove", "-pgpass=pass", "-debug"},
			expectFatal: true,
		},
		{
			name:        "importand db [template0] removal",
			args:        []string{"cmd", "-databases=template0", "-operation=remove", "-pgpass=pass", "-debug"},
			expectFatal: true,
		},
		{
			name:        "importand db [template1] removal",
			args:        []string{"cmd", "-databases=template1", "-operation=remove", "-pgpass=pass", "-debug"},
			expectFatal: true,
		},
		{
			name:        "many importand db removal",
			args:        []string{"cmd", "-databases=template1,postgres", "-operation=remove", "-pgpass=pass", "-debug"},
			expectFatal: true,
		},
		{
			name:        "only halper flag",
			args:        []string{"cmd", "-h"},
			expectFatal: false,
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

			// Тестируем логику сразу при возможном наличии вспомогательных флагов
			if args.HavingHelpFlag {
				ce.HelpFlagsExecution(args, testLogger)
				if args.OperationType == "none" { // если нет операции с СУБД, то нет смысла подключаться к БД, просто исполняем вспомогательные флаги и завершаем
					testLogger.Warn("ВНИМАНИЕ! Вы не ввели команду для работы с СУБД")
					return
				}
			}

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

			defer func() { // если ожидаемая ошибка, то это не считается фэйлом теста
				if r := recover(); r != nil {
					if !tt.expectFatal {
						t.Errorf("Неожиданная ошибка: %v", r)
					} else {
						t.Logf("Ожидаемая ошибка: %v", r)
					}
				}
			}()

			ctx := context.Background()
			ce.DBMSFlagsExecution(ctx, mock, args, testLogger)

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("Ожидаемые запросы не выполнены: %v", err)
			}
		})
	}
}
