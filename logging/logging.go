package logging

import (
	"os"

	"github.com/sirupsen/logrus"
)

var log = logrus.New()

func InitLog() {
	log.SetReportCaller(true)
	log.SetFormatter(&logrus.TextFormatter{
		DisableColors: true,
		FullTimestamp: true,
	})
	file, err := os.OpenFile("logs/logrus.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err == nil {
		log.Out = file
	} else {
		logrus.Info("Failed to log to file, using default stderr")
	}
}

func GetLogger() *logrus.Logger {
	return log
}
