package CommandProcessing

import (
	"flag"
	"os"
	"strings"

	"github.com/sirupsen/logrus"
)

func CommandProcessing(logger *logrus.Logger) CommandArgs {
	logger.Infof("Получили командную строку: %q", strings.Join(os.Args, " "))

	helpDump := flag.Bool("h", false, "Вывод опций")
	logger.Debugf("Заполнили helpDump по умолчанию: %v", *helpDump)

	debugInfo := flag.Bool("debug", false, "Включить дебаг информацию")
	logger.Debugf("Заполнили debugInfo по умолчанию: %v", *debugInfo)

	dbNamesWoParse := flag.String("databases", "none", "Имена баз данных в формате -databases=db1,db2,template%")
	logger.Debugf("Заполнили имена баз данных значением по умолчанию: %s", *dbNamesWoParse)

	operationType := flag.String("operation", "none", "Тип операции: backup|remove в формате -operation=")
	logger.Debugf("Заполнили тип операции значением по умолчанию: %s", *operationType)

	postgresPasswordURL := flag.String("pgpass", "none", "Пароль от вашего пользователя postgres на машине (если в вашем pg_hba.conf не установлен trust для local) в формате -pgpass=")
	logger.Debugf("Заполнили postgresPassword значением по умолчанию: %s", *postgresPasswordURL)

	flag.Parse()
	logger.Debugf("Считали Имена баз данных: %s", *dbNamesWoParse)
	logger.Debugf("Считали Тип Операции: %s", *operationType)
	logger.Debugf("Считали postgresPassword: %s", *postgresPasswordURL)
	logger.Debugf("Считали helpDump: %v", *helpDump)
	logger.Debugf("Считали debugInfo: %v", *debugInfo)

	if *operationType != "none" {
		if *dbNamesWoParse == "none" {
			logger.Fatal("Не введены имена баз данных. *Требуемый формат: -databases=db1,db2")
		}

		if *postgresPasswordURL == "none" {
			logger.Warn("ВНИМАНИЕ! Вы не ввели пароль от пользователя postgres, проверьте, что в вашем pg_hba.conf установлен trust для local")
		}

	} else if *helpDump {
	} else {
		logger.Fatal("Не введён тип операции")
	}

	dbNames := strings.Split(*dbNamesWoParse, ",")
	logger.Debugf("Распарсили имена строк в массив: %q (тип: %T)", dbNames, dbNames)

	logger.Info("Успешно обработали команду")

	return CommandArgs{DbNames: dbNames, OperationType: *operationType, PostgresPasswordURL: *postgresPasswordURL, HelperDump: *helpDump, DebugInfo: *debugInfo}
}
