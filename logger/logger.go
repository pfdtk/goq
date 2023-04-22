package logger

import "go.uber.org/zap"

type Logger interface {
	Debug(args ...interface{})
	Debugf(template string, args ...interface{})
	Info(args ...interface{})
	Infof(template string, args ...interface{})
	Warn(args ...interface{})
	Warnf(template string, args ...interface{})
	Error(args ...interface{})
	Errorf(template string, args ...interface{})
	Fatal(args ...interface{})
	Fatalf(template string, args ...interface{})
}

var logger Logger

func SetLogger(l Logger) {
	logger = l
}

func GetLogger() Logger {
	if logger == nil {
		return GetDefaultLogger()
	}
	return logger
}

func GetDefaultLogger() Logger {
	z, err := zap.NewDevelopment()
	if err != nil {
		return nil
	}
	return z.Sugar()
}
