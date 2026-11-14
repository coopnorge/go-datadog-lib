package tracing

import (
	"context"

	"github.com/coopnorge/go-logger"
)

// LogWithTrace will log message by logger.Level with trace if it's present in context.Context
//
// Deprecated: LogWithTrace will be removed in a future major version and
// should not be used. Use logger.WithContext(ctx Context) from
// github.com/coopnorge/go-logger.
// The whole module is deprecated; use the loghook in
// github.com/coopnorge/go-telemetry-lib/loghook, which adds the trace and span IDs to log
// entries.
func LogWithTrace(ctx context.Context, severity logger.Level, message string) {
	entry := logger.WithContext(ctx)
	logWithSeverity(entry, severity, message)
}

// LogFieldsWithTrace will log message by logger.Level with trace if it's present in context.Context
//
// Deprecated: LogFieldsWithTrace will be removed in a future major version and
// should not be used. Use logger.WithContext(ctx Context) from
// github.com/coopnorge/go-logger.
// The whole module is deprecated; use the loghook in
// github.com/coopnorge/go-telemetry-lib/loghook, which adds the trace and span IDs to log
// entries.
func LogFieldsWithTrace(ctx context.Context, severity logger.Level, message string, fields logger.Fields) {
	entry := logger.WithContext(ctx).WithFields(fields)
	logWithSeverity(entry, severity, message)
}

func logWithSeverity(entry *logger.Entry, severity logger.Level, message string) {
	switch severity {
	case logger.LevelFatal:
		entry.Fatal(message)
	case logger.LevelError:
		entry.Error(message)
	case logger.LevelWarn:
		entry.Warn(message)
	case logger.LevelInfo:
		entry.Info(message)
	case logger.LevelDebug:
		entry.Debug(message)
	}
}
