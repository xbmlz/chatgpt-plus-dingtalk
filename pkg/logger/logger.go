package logger

import (
	log "github.com/sirupsen/logrus"
)

var logger *log.Logger

func Initialize(level string) {
	logger = log.New()
	formatter := log.TextFormatter{
		ForceColors:               true,
		EnvironmentOverrideColors: true,
		TimestampFormat:           "2006-01-02 15:04:05",
		FullTimestamp:             true,
	}

	logger.SetFormatter(&formatter)
	if level == "debug" {
		logger.SetLevel(log.DebugLevel)
		logger.SetReportCaller(true)
	} else {
		logger.SetLevel(log.InfoLevel)
		logger.SetReportCaller(false)
	}
}

func Debug(args ...interface{}) {
	logger.Debug(args...)
}

func Info(args ...interface{}) {
	logger.Info(args...)
}

func Warn(args ...interface{}) {
	logger.Warn(args...)
}

func Error(args ...interface{}) {
	logger.Error(args...)
}

func Fatal(args ...interface{}) {
	logger.Fatal(args...)
}
