// Package tracelogger provides a go-logger hook that adds Datadog trace
// context to log entries.
//
// Deprecated: The whole module is deprecated. Use NewHook from
// github.com/coopnorge/go-telemetry-lib/loghook. See the root package
// documentation for the migration guide and ADR.
package tracelogger

import (
	"github.com/DataDog/dd-trace-go/v2/ddtrace/tracer"
	"github.com/coopnorge/go-logger"
)

// DDContextLogHook ensures that any span in the log context is correlated to log output.
//
// Deprecated: The whole module is deprecated. Use NewHook from
// github.com/coopnorge/go-telemetry-lib/loghook.
type DDContextLogHook struct{}

// Fire implements logger.Hook interface, attaches trace and span details found in entry context
//
// Deprecated: The whole module is deprecated. Use NewHook from
// github.com/coopnorge/go-telemetry-lib/loghook.
func (d *DDContextLogHook) Fire(he *logger.HookEntry) (bool, error) {
	ctx := he.Context
	if ctx == nil {
		return false, nil
	}
	span, found := tracer.SpanFromContext(he.Context)
	if !found {
		return false, nil
	}
	he.Data["dd.trace_id"] = span.Context().TraceID()
	he.Data["dd.span_id"] = span.Context().SpanID()
	return true, nil
}

// NewHook will create a new Hook compatible with go-logger,
// to automatically extract Span/Trace information from the Log-entry's context.Context,
// and add them as fields to the log-message.
//
// Deprecated: The whole module is deprecated. Use NewHook from
// github.com/coopnorge/go-telemetry-lib/loghook.
func NewHook() logger.Hook {
	return &DDContextLogHook{}
}
