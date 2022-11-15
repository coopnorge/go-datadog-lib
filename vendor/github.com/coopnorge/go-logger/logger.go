package logger

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/sirupsen/logrus"
)

// Fields type, used to pass to `WithFields`.
type Fields map[string]interface{}

type NowFunc func() time.Time

// Logger is our logger with the needed structured logger we use
type Logger struct {
	logrusLogger *logrus.Logger
	now          NowFunc
	output       io.Writer
	level        Level
	reportCaller bool
}

func (logger *Logger) applyOptions(opts ...LoggerOption) {
	for _, opt := range opts {
		opt.Apply(logger)
	}
	logger.logrusLogger.SetOutput(logger.output)
	logger.logrusLogger.SetLevel(mapLevelToLogrusLevel(logger.level))
}

func New(opts ...LoggerOption) *Logger {
	logrusLogger := logrus.New()
	logrusLogger.SetFormatter(&logrus.JSONFormatter{})
	logger := &Logger{
		logrusLogger: logrusLogger,
		now:          NowFunc(time.Now),
		output:       os.Stdout,
		level:        LevelWarn,
		reportCaller: true,
	}
	logger.applyOptions(opts...)
	return logger
}

func (logger *Logger) entry() Entry {
	fields := logrus.Fields{}
	if logger.reportCaller {
		frame := getCaller()
		fields["file"] = fmt.Sprintf("%s:%v", frame.File, frame.Line)
		fields["function"] = frame.Function
	}
	return logger.logrusLogger.WithTime(logger.now()).WithFields(fields)
}

// WithFields forwards a logging call with fields
func (logger *Logger) WithFields(fields Fields) Entry {
	return logger.logrusLogger.WithTime(logger.now()).WithFields(logrus.Fields(fields))
}

// Infof forwards a logging call in the (format, args) format
func (logger *Logger) Info(args ...interface{}) {
	logger.entry().Info(args...)
}

// Infof forwards a logging call in the (format, args) format
func (logger *Logger) Infof(format string, args ...interface{}) {
	logger.entry().Infof(format, args...)
}

// Error forwards an error logging call
func (logger *Logger) Error(args ...interface{}) {
	logger.entry().Error(args...)
}

// Errorf forwards an error logging call
func (logger *Logger) Errorf(format string, args ...interface{}) {
	logger.entry().Errorf(format, args...)
}

// Debug forwards a debugging logging call
func (logger *Logger) Debug(args ...interface{}) {
	logger.entry().Debug(args...)
}

// Debugf forwards a debugging logging call
func (logger *Logger) Debugf(format string, args ...interface{}) {
	logger.entry().Debugf(format, args...)
}

// Warn forwards a warning logging call
func (logger *Logger) Warn(args ...interface{}) {
	logger.entry().Warn(args...)
}

// Warnf forwards a warning logging call
func (logger *Logger) Warnf(format string, args ...interface{}) {
	logger.entry().Warnf(format, args...)
}

// Fatal forwards a fatal logging call
func (logger *Logger) Fatal(args ...interface{}) {
	logger.entry().Fatal(args...)
}

// Fatalf forwards a fatal logging call
func (logger *Logger) Fatalf(format string, args ...interface{}) {
	logger.entry().Fatalf(format, args...)
}
