package echo

import (
	"encoding/json"
	"github.com/coopnorge/go-logger"
	"io"

	"github.com/labstack/gommon/log"
)

// WrappedEchoLogger that can be passed to Echo middleware for Datadog integration
// implements Echo Logger from vendor/github.com/labstack/echo/v4/log.go
type WrappedEchoLogger struct {
	log    logger.Entry
	output io.Writer

	// prefix for logs
	prefix string
	level  logger.Level
}

// NewWrappedEchoLogger instance
func NewWrappedEchoLogger() *WrappedEchoLogger {
	return &WrappedEchoLogger{
		log:    logger.New(logger.WithLevel(logger.LevelInfo)),
		prefix: "Datadog Echo",
		level:  logger.LevelInfo,
	}
}

// Output not supported to access in Coop logger, so it's returns only stub
func (wel *WrappedEchoLogger) Output() io.Writer {
	return wel.output
}

// SetOutput not supported to change in Coop logger, accepting only stub
func (wel *WrappedEchoLogger) SetOutput(w io.Writer) {
	wel.output = w
}

func (wel *WrappedEchoLogger) Prefix() string {
	return wel.prefix
}

func (wel *WrappedEchoLogger) SetPrefix(p string) {
	wel.prefix = p
}

func (wel *WrappedEchoLogger) Level() log.Lvl {
	switch wel.level {
	case logger.LevelDebug:
		return log.DEBUG
	case logger.LevelInfo:
		return log.INFO
	case logger.LevelWarn:
		return log.WARN
	case logger.LevelError:
		return log.ERROR
	case logger.LevelFatal:
		return log.ERROR
	default:
		return log.OFF
	}
}

func (wel *WrappedEchoLogger) SetLevel(v log.Lvl) {
	switch v {
	case log.DEBUG:
		wel.level = logger.LevelDebug
	case log.INFO:
		wel.level = logger.LevelInfo
	case log.WARN:
		wel.level = logger.LevelWarn
	case log.ERROR:
		wel.level = logger.LevelError
	case log.OFF:
		return // Ignore Coop logger cannot be disabled yet
	}

	logger.SetLevel(wel.level)
}

// SetHeader not supported
func (wel *WrappedEchoLogger) SetHeader(_ string) {
	return
}

func (wel *WrappedEchoLogger) Print(i ...interface{}) {
	wel.log.Info(i...)
}

func (wel *WrappedEchoLogger) Printf(format string, args ...interface{}) {
	wel.log.Infof(format, args...)
}

func (wel *WrappedEchoLogger) Printj(j log.JSON) {
	wel.log.Info(wel.jsonToString(j))
}

func (wel *WrappedEchoLogger) Debug(i ...interface{}) {
	wel.log.Debug(i...)
}

func (wel *WrappedEchoLogger) Debugf(format string, args ...interface{}) {
	wel.log.Debugf(format, args...)
}

func (wel *WrappedEchoLogger) Debugj(j log.JSON) {
	wel.log.Debug(wel.jsonToString(j))
}

func (wel *WrappedEchoLogger) Info(i ...interface{}) {
	wel.log.Info(i...)
}

func (wel *WrappedEchoLogger) Infof(format string, args ...interface{}) {
	wel.log.Infof(format, args...)
}

func (wel *WrappedEchoLogger) Infoj(j log.JSON) {
	wel.log.Info(wel.jsonToString(j))
}

func (wel *WrappedEchoLogger) Warn(i ...interface{}) {
	wel.log.Warn(i...)
}

func (wel *WrappedEchoLogger) Warnf(format string, args ...interface{}) {
	wel.log.Warnf(format, args...)
}

func (wel *WrappedEchoLogger) Warnj(j log.JSON) {
	wel.log.Warn(wel.jsonToString(j))
}

func (wel *WrappedEchoLogger) Error(i ...interface{}) {
	wel.log.Error(i...)
}

func (wel *WrappedEchoLogger) Errorf(format string, args ...interface{}) {
	wel.log.Errorf(format, args...)
}

func (wel *WrappedEchoLogger) Errorj(j log.JSON) {
	wel.log.Error(wel.jsonToString(j))
}

func (wel *WrappedEchoLogger) Fatal(i ...interface{}) {
	wel.log.Fatal(i...)
}

func (wel *WrappedEchoLogger) Fatalj(j log.JSON) {
	wel.log.Fatal(wel.jsonToString(j))
}

func (wel *WrappedEchoLogger) Fatalf(format string, args ...interface{}) {
	wel.log.Fatalf(format, args...)
}

func (wel *WrappedEchoLogger) Panic(i ...interface{}) {
	wel.log.Fatal(i...)
}

func (wel *WrappedEchoLogger) Panicj(j log.JSON) {
	wel.log.Fatal(wel.jsonToString(j))
}

func (wel *WrappedEchoLogger) Panicf(format string, args ...interface{}) {
	wel.log.Fatalf(format, args...)
}

func (wel *WrappedEchoLogger) jsonToString(j log.JSON) string {
	jb, marshalErr := json.Marshal(j)
	if marshalErr != nil {
		wel.log.Errorf("unable to marshal json data for logs, marshal error: :%v", marshalErr)
		return ""
	}

	return string(jb)
}
