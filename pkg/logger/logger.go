package logger

import (
	log "github.com/sirupsen/logrus"
)

var Logger *log.Logger

func Init(level string) {
	Logger = log.New()
	if level == "debug" {
		Logger.SetLevel(log.DebugLevel)
	} else {
		Logger.SetLevel(log.InfoLevel)
	}
}

func Info(args ...interface{}) {
	Logger.Info(args)
}

func Warning(args ...interface{}) {
	Logger.Warn(args)
}

func Debug(args ...interface{}) {
	Logger.Debug(args)
}

func Error(args ...interface{}) {
	Logger.Error(args)
}

func Fatal(args ...interface{}) {
	Logger.Fatal(args)
}
