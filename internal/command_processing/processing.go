package CommandProcessing

import (
	"flag"
	"strings"

	"github.com/sirupsen/logrus"
)

func CommandProcessing(logger *logrus.Logger) CommandArgs {
	dbNamesWoParse := flag.String("databases", "none", "Имена баз данных в формате -databses=db1,db2")
	logger.Debugf("Заполнили имена баз данных значением по умолчанию: %s", *dbNamesWoParse)

	operationType := flag.String("operation", "none", "Тип операции: backup|remove")
	logger.Debugf("Заполнили тип операции значением по умолчанию: %s", *operationType)

	postgresPasswordURL := flag.String("pgpass", "none", "Пароль от вашего пользователя postgres на машине")
	logger.Debugf("Заполнили postgresPassword значением по умолчанию: %s", *postgresPasswordURL)

	flag.Parse()
	logger.Debugf("Считали Имена баз данных: %s", *dbNamesWoParse)
	logger.Debugf("Считали Тип Операции: %s", *operationType)
	logger.Debugf("Считали postgresPassword: %s", *postgresPasswordURL)

	if *dbNamesWoParse == "none" {
		logger.Fatal("Не введены имена баз данных. *Требуемый формат: -databses=db1,db2")
	}

	if *operationType == "none" {
		logger.Fatal("Не введён тип операции")
	}

	if *postgresPasswordURL == "none" {
		logger.Fatal("Не введён postgresPassword")
	}

	dbNames := strings.Split(*dbNamesWoParse, ",")
	logger.Debugf("Распарсили имена строк в массив: %v (тип: %T)", dbNames, dbNames)

	logger.Info("Успешно обработали команду")

	return CommandArgs{DbNames: dbNames, OperationType: *operationType, PostgresPasswordURL: *postgresPasswordURL}
}
