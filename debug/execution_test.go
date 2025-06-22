package UnitTests

import (
	"context"
	"flag"
	"os"
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
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil {
					t.Errorf("Неожиданная ошибка: %v", r)
				}
			}()

			flag.CommandLine = flag.NewFlagSet("cmd", flag.ExitOnError) // Сбрасываем для кажого теста
			os.Args = tt.args                                           // Переопределяем аргументы коммандной строки для конкретного теста

			args := cp.CommandProcessing(testLogger)

			mock, err := pgxmock.NewConn()
			if err != nil {
				t.Fatalf("Не удалось создать mock подключение: %v", err)
			}
			defer mock.Close(context.Background())

			for _, dbName := range args.DbNames {
				mock.ExpectExec("DROP DATABASE IF EXISTS " + dbName).
					WillReturnResult(pgxmock.NewResult("DROP", 1))
			}

			ctx := context.Background()
			ce.CommandExecution(ctx, mock, args, testLogger)

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("Ожидаемые запросы не выполнены: %v", err)
			}
		})
	}
}
