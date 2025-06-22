package CommandExecution

import (
	"context"
	"fmt"

	cp "github.com/fl4cko11/dbc-utility/internal/command_processing"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/sirupsen/logrus"
)

type DBConn interface { // Для возможности mock-тестирования
	Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error) // ввод данного интерфейса обоснован тем, что утилита только управляет СУБД (без запросов и тп)
}

func CommandExecution(ctx context.Context, conn DBConn, args cp.CommandArgs, logger *logrus.Logger) {
	if args.OperationType == "remove" {
		for _, val := range args.DbNames {
			logger.Debugf("Приняли имя БД для удаления: %q", val)

			_, erre := conn.Exec(ctx, fmt.Sprintf("DROP DATABASE %s", val))
			if erre != nil {
				logger.Errorf("Не удалось удалить базу данных %s", val)
			} else {
				logger.Infof("Успешно удалили БД %s", val)
			}
		}
		logger.Info("Завершили процесс удаления заданных БД")
	} // else if args.OperationType == "backup" {}
}
