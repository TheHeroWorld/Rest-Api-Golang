package logging

import (
	"io"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/sirupsen/logrus/hooks/writer"
)

type Hook struct {
	Writer    io.Writer
	LogLevels []logrus.Level
}

// Вызов базвого хука который реагирует на level Panic, fatal, Error Warn
func InitHook() {
	log.AddHook(&writer.Hook{
		Writer: os.Stdout,
		LogLevels: []logrus.Level{
			logrus.PanicLevel,
			logrus.FatalLevel,
			logrus.ErrorLevel,
			logrus.WarnLevel},
	})
}

// Хук который ловит все записи с указаными левелами и выводит ошикби в консоль
func (hook *Hook) Fire(entry *logrus.Entry) error {
	line, err := entry.Bytes()
	if err != nil {
		return err
	}
	_, err = hook.Writer.Write(line)
	return err
}

// Метод Levels который определяет уровни для подвхата
func (hook *Hook) Levels() []logrus.Level {
	return hook.LogLevels
}
