package logs

import (
	"io"

	"github.com/sirupsen/logrus"
)

type DynamicFormatter struct {
	InfoFormatter  logrus.Formatter
	ErrorFormatter logrus.Formatter
}

func (f *DynamicFormatter) Format(entry *logrus.Entry) ([]byte, error) { // модифицировал логер, чтобы на инфо уровне не было строк и пути
	switch entry.Level {
	case logrus.InfoLevel:
		dataCopy := make(logrus.Fields, len(entry.Data))
		for k, v := range entry.Data {
			dataCopy[k] = v
		}
		delete(dataCopy, "caller")

		entryCopy := *entry
		entryCopy.Data = dataCopy

		return f.InfoFormatter.Format(&entryCopy)
	default:
		return f.ErrorFormatter.Format(entry)
	}
}

func InitLogger(out io.Writer) *logrus.Logger {
	logger := logrus.New()

	logger.SetReportCaller(true)
	logger.SetLevel(logrus.TraceLevel)

	logger.SetFormatter(&DynamicFormatter{
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
