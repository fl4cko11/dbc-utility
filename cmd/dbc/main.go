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
	// пользователю не нужна debug информация исполнения по умолчанию
	// тем более информация о процессе парсинга
	// она доступна только в тестах на уровне разработчика
	// если хотим увидеть, то можно реализовать на уровне main парсинг флага -debug
	logger := logs.InitLogger(os.Stderr, false)

	args := cp.CommandProcessing(logger)

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
	logger.Info("Успешно установили соединение с БД")

	ce.CommandExecution(ctx, conn, args, logger)
}
