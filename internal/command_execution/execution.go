package CommandExecution

import (
	"context"
	"fmt"
	"net/url"

	cp "github.com/fl4cko11/dbc-utility/internal/command_processing"
	"github.com/jackc/pgx/v5"
	"github.com/sirupsen/logrus"
)

func CommandExecution(logger *logrus.Logger) {
	args := cp.CommandProcessing(logger)

	dbURL := &url.URL{
		Scheme: "postgres",
		Host:   "localhost",
		User:   url.UserPassword("postgres", args.PostgresPasswordURL),
		Path:   "/postgres",
	}
	dbURLs := dbURL.String()
	logger.Debugf("Конвертировали dbURL в строку: %v", dbURLs)

	ctx := context.Background()
	conn, errc := pgx.Connect(ctx, dbURLs)
	if errc != nil {
		logger.Fatal("Не удалось подключиться к базе данных")
	}
	defer conn.Close(ctx)
	logger.Info("Успешно установили соединение с БД")

	if args.OperationType == "remove" {
		for _, val := range args.DbNames {
			logger.Debugf("Приняли имя БД для удаления: %v", val)

			_, erre := conn.Exec(ctx, fmt.Sprintf("DROP DATABASE IF EXISTS %s", val))
			if erre != nil {
				logger.Errorf("Не удалось удалить базу данных %v", val)
			} else {
				logger.Infof("Успешно удалили БД %v", val)
			}
		}
		logger.Info("Успешно удалили заданные БД")
	} else if args.OperationType == "backup" {

	}
}
