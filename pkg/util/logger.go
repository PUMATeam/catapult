package util

import (
	"os"

	logrus "github.com/sirupsen/logrus"
)

// InitLog initializes the logger
func InitLog() *logrus.Logger {
	// TODO make configurable
	logger := logrus.New()
	logger.SetOutput(os.Stdout)
	logger.SetLevel(logrus.DebugLevel)
	logger.Formatter = &logrus.TextFormatter{
		FullTimestamp: true,
	}
	return logger
}
