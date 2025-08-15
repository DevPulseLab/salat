package service

import (
	"io"
	"os"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

type ErrorLogger struct {
	logger    *logrus.Logger
	initOnce  sync.Once
	toConsole bool
	path      string
	level     logrus.Level
}

func NewErrorLogger(path string, toConsole bool, level logrus.Level) *ErrorLogger {
	return &ErrorLogger{
		logger:    logrus.New(),
		toConsole: toConsole,
		path:      path,
		level:     level,
	}
}

func (el *ErrorLogger) GetDefaultLogger() *logrus.Logger {
	el.initOnce.Do(func() {
		lumberjackLogger := &lumberjack.Logger{
			Filename:   el.dailyLogFileName(),
			MaxSize:    10,
			MaxAge:     10,
			MaxBackups: 0,
			Compress:   true,
		}

		var output io.Writer
		if el.toConsole {
			output = io.MultiWriter(os.Stdout, lumberjackLogger)
		} else {
			output = lumberjackLogger
		}

		el.logger.SetOutput(output)
		el.logger.SetLevel(el.level)
		el.logger.SetFormatter(&logrus.JSONFormatter{
			TimestampFormat: time.RFC3339,
			PrettyPrint:     false,
		})
	})
	return el.logger
}

func (el *ErrorLogger) dailyLogFileName() string {
	today := time.Now().Format("2006-01-02")
	return el.path + "_" + today + ".log"
}
