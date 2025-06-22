package logs

import (
	"io"

	"github.com/sirupsen/logrus"
)

type dynamicFormatter struct {
	InfoFormatter  logrus.Formatter
	ErrorFormatter logrus.Formatter
}

func (f *dynamicFormatter) Format(entry *logrus.Entry) ([]byte, error) { // модифицировал логер, чтобы на инфо уровне не было строк и пути
	if entry.Level == logrus.InfoLevel || entry.Level == logrus.WarnLevel {
		newEntry := &logrus.Entry{
			Logger:  entry.Logger,
			Data:    make(logrus.Fields),
			Time:    entry.Time,
			Level:   entry.Level,
			Message: entry.Message,
		}

		for k, v := range entry.Data {
			if k != "caller" {
				newEntry.Data[k] = v
			}
		}

		return f.InfoFormatter.Format(newEntry)
	}
	return f.ErrorFormatter.Format(entry)
}

func InitLogger(out io.Writer, debugInfo bool) *logrus.Logger {
	logger := logrus.New()

	logger.SetReportCaller(true)

	if debugInfo {
		logger.SetLevel(logrus.TraceLevel)
	} else {
		logger.SetLevel(logrus.InfoLevel)
	}

	logger.SetFormatter(&dynamicFormatter{
		InfoFormatter: &logrus.TextFormatter{
			TimestampFormat: "2006-01-02 15:04:05",
			ForceColors:     true,
		},
		ErrorFormatter: &logrus.TextFormatter{
			TimestampFormat: "2006-01-02 15:04:05",
			ForceColors:     true,
		},
	})

	return logger
}
