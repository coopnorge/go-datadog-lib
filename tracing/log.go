package tracing

import (
	"context"
	"fmt"

	"github.com/coopnorge/go-logger"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

// LogWithTrace will log message by logger.Level with trace if it's present in context.Context
//
// Deprecated: LogWithTrace will be removed in a future major version and
// should not be used. Use logger.WithContext(ctx Context) from
// github.com/coopnorge/go-logger
func LogWithTrace(sourceCtx context.Context, severity logger.Level, message string) {
	messageToLog := getMessageToLog(sourceCtx, message)
	emptyEntry := logger.WithFields(map[string]interface{}{})

	logWithSeverity(emptyEntry, severity, messageToLog)
}

// LogFieldsWithTrace will log message by logger.Level with trace if it's present in context.Context
//
// Deprecated: LogFieldsWithTrace will be removed in a future major version and
// should not be used. Use logger.WithContext(ctx Context) from
// github.com/coopnorge/go-logger
func LogFieldsWithTrace(sourceCtx context.Context, severity logger.Level, message string, fields logger.Fields) {
	messageToLog := getMessageToLog(sourceCtx, message)
	entry := logger.WithFields(fields)

	logWithSeverity(entry, severity, messageToLog)
}

func getMessageToLog(ctx context.Context, message string) string {
	var messageToLog string

	span, exists := tracer.SpanFromContext(ctx)
	if exists {
		messageToLog = fmt.Sprintf("%s %v dd.lang=go", message, span)
	} else {
		messageToLog = message
	}

	return messageToLog
}

func logWithSeverity(entry logger.Entry, severity logger.Level, message string) {
	switch severity {
	case logger.LevelFatal:
		entry.Fatalf(message)
	case logger.LevelError:
		entry.Errorf(message)
	case logger.LevelWarn:
		entry.Warnf(message)
	case logger.LevelInfo:
		entry.Infof(message)
	case logger.LevelDebug:
		entry.Debugf(message)
	}
}
