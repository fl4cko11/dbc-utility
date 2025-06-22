package CommandExecution

import (
	"context"
	"fmt"
	"strings"

	cp "github.com/fl4cko11/dbc-utility/internal/command_processing"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/sirupsen/logrus"
)

type DBConn interface { // Для возможности mock-тестирования
	Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
}

func CommandExecution(ctx context.Context, conn DBConn, args cp.CommandArgs, logger *logrus.Logger) {
	if args.OperationType == "remove" {
		for _, val := range args.DbNames {
			logger.Debugf("Приняли имя БД для удаления: %q", val)

			if strings.EqualFold(val, "postgres") || strings.EqualFold(val, "template0") || strings.EqualFold(val, "template1") {
				logger.Fatalf("Попытка удалить обязательную бд: %s", val)
			}

			if strings.Contains(val, "%") {
				logger.Debugf("Начали процесс удаления по шаблонному имени %s", val)

				rows, errQ := conn.Query(ctx, fmt.Sprintf("SELECT datname FROM pg_database WHERE datname LIKE '%s' AND datistemplate = false AND datname NOT IN ('postgres','template0','template1');", val))
				if errQ != nil {
					logger.Fatalf("Не удалось получить имена БД из СУБД по шаблону %s, ошибка: %v", val, errQ)
				}
				defer rows.Close()
				logger.Debug("Получили строки с именами БД для шаблона")

				var dbName string
				var dbList []string
				for rows.Next() {
					errS := rows.Scan(&dbName)
					if errS != nil {
						logger.Fatalf("Не удалось считать полученную строку из запроса к именам БД, ошибка: %v", errS)
					}
					logger.Debugf("Считали имя БД соответствующее шаблону %s: %s", val, dbName)
					dbList = append(dbList, dbName)
				}

				logger.Debugf("Начинаем удаление всех бд для шаблона %s", val)
				for _, dbname := range dbList {
					_, errE := conn.Exec(ctx, fmt.Sprintf("DROP DATABASE %s", dbname))
					if errE != nil {
						logger.Errorf("Не удалось удалить БД %s, ошибка: %v", dbname, errE)
					} else {
						logger.Infof("Успешно удалили БД соответствующее шаблону %s: %s", val, dbname)
					}
				}
			} else {
				_, erre := conn.Exec(ctx, fmt.Sprintf("DROP DATABASE %s", val))
				if erre != nil {
					logger.Errorf("Не удалось удалить БД %s, ошибка: %v", val, erre)
				} else {
					logger.Infof("Успешно удалили БД %s", val)
				}
			}
		}
		logger.Info("Завершили процесс удаления заданных БД")
	} // else if args.OperationType == "backup" {}
}
