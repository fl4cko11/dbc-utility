package CommandProcessing

import (
	"flag"
	"os"
	"strings"

	"github.com/sirupsen/logrus"
)

func CommandProcessing(logger *logrus.Logger) CommandArgs {
	logger.Infof("Получили командную строку: %q", strings.Join(os.Args, " "))

	helpDump := flag.Bool("h", false, "Вывод опций флаг в формате -h")

	debugInfo := flag.Bool("debug", false, "Включить дебаг информацию флаг в формате -debug")

	dbNamesWoParse := flag.String("databases", "none", "Имена баз данных флаг в формате -databases=db1,db2,template%")

	operationType := flag.String("operation", "none", "Тип операции: backup|remove флаг в формате -operation=")

	postgresPasswordURL := flag.String("pgpass", "none", "Пароль от вашего пользователя postgres на машине (если в вашем pg_hba.conf не установлен trust для local) флаг в формате -pgpass=")

	flag.Parse()
	logger.Infof("Считали Имена баз данных: %s", *dbNamesWoParse)
	logger.Infof("Считали Тип Операции: %s", *operationType)
	logger.Infof("Считали postgresPassword: %s", *postgresPasswordURL)
	logger.Infof("Считали helpDump: %v", *helpDump)
	logger.Infof("Считали debugInfo: %v", *debugInfo)

	if *debugInfo { // если пользователь потребовал debug информацию
		logger.SetLevel(logrus.TraceLevel)
	} else {
		logger.SetLevel(logrus.InfoLevel)
	}

	args := CommandArgs{DbNames: []string{"none"}, OperationType: *operationType, HavingHelpFlag: false, HelperDump: false}

	if *operationType != "none" {
		if *dbNamesWoParse == "none" {
			logger.Fatal("Не введены имена баз данных. *Требуемый формат: -databases=db1,db2")
		}

		if *postgresPasswordURL == "none" {
			logger.Warn("ВНИМАНИЕ! Вы не ввели пароль от пользователя postgres, проверьте, что в вашем pg_hba.conf установлен trust для local")
		}

		dbNames := strings.Split(*dbNamesWoParse, ",")
		logger.Debugf("Распарсили имена строк в массив: %q (тип: %T)", dbNames, dbNames)
		args.DbNames = dbNames
	}
	logger.Info("Успешно обработали команду")

	args.HelperDump = *helpDump
	if *helpDump {
		args.HavingHelpFlag = true
	}

	return args
}
