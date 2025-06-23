package main

import (
	"context"
	"net/url"
	"os"

	ce "github.com/fl4cko11/dbc-utility/internal/command_execution"
	cp "github.com/fl4cko11/dbc-utility/internal/command_processing"
	"github.com/fl4cko11/dbc-utility/logs"
	"github.com/jackc/pgx/v5"
)

func main() {
	logger := logs.InitLogger(os.Stderr)

	args := cp.CommandProcessing(logger)

	if args.HavingHelpFlag {
		ce.HelpFlagsExecution(args, logger)
		if args.OperationType == "none" { // если нет операции с СУБД, то нет смысла подключаться к БД, просто исполняем вспомогательные флаги и завершаем
			logger.Warn("ВНИМАНИЕ! Вы не ввели команду для работы с СУБД")
			return
		}
	}

	dbURL := &url.URL{
		Scheme: "postgres",
		Host:   "localhost",
		User:   url.UserPassword("postgres", args.PostgresPasswordURL),
		Path:   "/postgres",
	}
	dbURLs := dbURL.String()
	logger.Debugf("Конвертировали dbURL в строку: %s", dbURLs)

	ctx := context.Background()
	conn, errc := pgx.Connect(ctx, dbURLs)
	if errc != nil {
		logger.Fatalf("Не удалось подключиться к базе данных, ошибка: %v", errc)
	}
	defer conn.Close(ctx)
	logger.Info("Успешно установили соединение с СУБД")

	ce.DBMSFlagsExecution(ctx, conn, args, logger)
}
