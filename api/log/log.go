package log

import (
	"go.uber.org/zap"
)

var GlobalLogger = NewLogger()

// NewLogger create a new SugaredLogger
func NewLogger() *zap.SugaredLogger {
	logger, _ := zap.NewProduction()
	return logger.Sugar()
}

// Infof uses fmt.Sprintf to log a templated message.
func Infof(template string, args ...interface{}) {
	GlobalLogger.Infof(template, args...)
}

// Warnf uses fmt.Sprintf to log a templated message.
func Warnf(template string, args ...interface{}) {
	GlobalLogger.Warnf(template, args...)
}

// Errorf uses fmt.Sprintf to log a templated message.
func Errorf(template string, args ...interface{}) {
	GlobalLogger.Errorf(template, args...)
}

// Debugf uses fmt.Sprintf to log a templated message.
func Debugf(template string, args ...interface{}) {
	GlobalLogger.Debugf(template, args...)
}

// Fatalf uses fmt.Sprintf to log a templated message.
func Fatalf(template string, args ...interface{}) {
	GlobalLogger.Fatalf(template, args...)
}

// Panicf uses fmt.Sprintf to log a templated message.
func Panicf(template string, args ...interface{}) {
	GlobalLogger.Panicf(template, args...)
}
