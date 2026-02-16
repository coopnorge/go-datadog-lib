package log // nolint:revive

import (
	"errors"
	"strings"

	"github.com/coopnorge/go-logger"
)

// Logger is a logging adapter between github.com/DataDog/dd-trace-go/v2
// an go-logger, do not create this directly, use NewLogger()
type Logger struct {
	instance *logger.Logger
}

// NewLogger creates a new Datadog logger that passes message to go-logger
//
// To inject the logger into tracer use
//
//	package main
//
//	import (
//		"github.com/coopnorge/go-logger/adapter/datadog"
//
//		"github.com/DataDog/dd-trace-go/v2/ddtrace/tracer"
//	)
//
//	func main() {
//		l, err := datadog.NewLogger(datadog.WithGlobalLogger())
//		tracer.UseLogger(l)
//	}
func NewLogger(opts ...LoggerOption) (*Logger, error) {
	logger := &Logger{}
	for _, opt := range opts {
		opt.Apply(logger)
	}
	if logger.instance == nil {
		return nil, errors.New("no go-logger instance provided, use WithGlobalLogger() or WithLogger() to configure the logger")
	}
	return logger, nil
}

// Log writes statements to the log
func (l *Logger) Log(msg string) {
	// Logs from github.com/DataDog/dd-trace-go/v2 will contain keywords
	// specifying the level of the log.
	if strings.Contains(msg, "ERROR") {
		l.instance.Error(msg)
		return
	}
	if strings.Contains(msg, "WARN") {
		l.instance.Warn(msg)
		return
	}
	if strings.Contains(msg, "INFO") {
		l.instance.Info(msg)
		return
	}
	if strings.Contains(msg, "DEBUG") {
		l.instance.Debug(msg)
		return
	}
	l.instance.WithField("datadog", "Datadog logger adapter could not resolve the logging level, falling back to warning").Warn(msg)
}
