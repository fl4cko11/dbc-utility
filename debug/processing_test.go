package UnitTests

import (
	"flag"
	"os"
	"testing"

	cp "github.com/fl4cko11/dbc-utility/internal/command_processing"
	"github.com/fl4cko11/dbc-utility/logs"
)

func TestCommandProcessing(t *testing.T) {
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
			name:        "valid args with single db",
			args:        []string{"cmd", "-databases=db1", "-operation=backup", "-pgpass=pass"},
			expectFatal: false,
		},
		{
			name:        "valid args with multiple dbs",
			args:        []string{"cmd", "-databases=db1,db2,db3", "-operation=remove", "-pgpass=secret"},
			expectFatal: false,
		},
		{
			name:        "valid args with single template",
			args:        []string{"cmd", "-databases=db1%", "-operation=backup", "-pgpass=pass"},
			expectFatal: false,
		},
		{
			name:        "valid args with args and templates",
			args:        []string{"cmd", "-databases=db1,db2%,db3%", "-operation=backup", "-pgpass=pass"},
			expectFatal: false,
		},
		{
			name:        "missing databases",
			args:        []string{"cmd", "-operation=backup", "-pgpass=pass"},
			expectFatal: true,
		},
		{
			name:        "missing operation",
			args:        []string{"cmd", "-databases=db1", "-pgpass=pass"},
			expectFatal: true,
		},
		{
			name:        "none values",
			args:        []string{"cmd", "-databases=none", "-operation=none", "-pgpass=none"},
			expectFatal: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			defer func() {
				if r := recover(); r != nil {
					if !tt.expectFatal {
						t.Errorf("Неожиданная ошибка: %v", r)
					} else {
						t.Logf("Ожидаемая ошибка: %v", r)
					}
				}
			}()

			flag.CommandLine = flag.NewFlagSet("cmd", flag.ExitOnError) // Сбрасываем для кажого теста
			os.Args = tt.args                                           // Переопределяем аргументы коммандной строки для конкретного теста

			cp.CommandProcessing(testLogger)
		})
	}
}
