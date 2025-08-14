package service

import (
	"io"
	"os"
	"sync"

	"github.com/sirupsen/logrus"
)

type ErrorLogger struct {
	log       *logrus.Logger
	file      *os.File
	initOnce  sync.Once
	initErr   error
	toConsole bool
	path      string
	level     logrus.Level
}

func NewErrorLogger(path string, toConsole bool, level logrus.Level) *ErrorLogger {
	return &ErrorLogger{
		log:       logrus.New(),
		toConsole: toConsole,
		path:      path,
		level:     level,
	}
}

func (el *ErrorLogger) GetDefaultLogger() (*logrus.Logger, error) {
	el.initOnce.Do(func() {
		file, err := os.OpenFile(el.path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
		if err != nil {
			el.initErr = err
			el.log.SetOutput(os.Stdout)
			return
		}
		el.file = file

		if el.toConsole {
			el.log.SetOutput(io.MultiWriter(os.Stdout, file))
		} else {
			el.log.SetOutput(file)
		}

		el.log.SetLevel(el.level)
		el.log.SetFormatter(&logrus.TextFormatter{FullTimestamp: true})
	})
	return el.log, el.initErr
}

func (el *ErrorLogger) Close() error {
	if el.file != nil {
		return el.file.Close()
	}
	return nil
}
