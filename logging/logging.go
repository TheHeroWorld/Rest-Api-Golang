package logging

import (
	"os"

	"github.com/sirupsen/logrus"
)

var log = logrus.New()

func InitLog() {
	//Запускаем ХУК
	InitHook()
	//Ниже настройки логгера
	log.SetReportCaller(true)
	log.SetFormatter(&logrus.TextFormatter{
		// Не работают цветной вывод, почему не знаю
		DisableColors: false,
		FullTimestamp: true,
	})
	file, err := os.OpenFile("logs/logrus.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err == nil {
		log.Out = file
	} else {
		logrus.Info("Failed to log to file, using default stderr")
	}
}

// Функция для вызова логгера в других пакетах
func GetLogger() *logrus.Logger {
	return log
}
