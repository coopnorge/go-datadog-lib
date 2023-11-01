package tracing

import "gopkg.in/DataDog/dd-trace-go.v1/ddtrace"

type (
	noopSpanContext struct{}
)

var _ ddtrace.SpanContext = (*noopSpanContext)(nil)

// ForeachBaggageItem implements ddtrace.SpanContext.
func (*noopSpanContext) ForeachBaggageItem(_ func(k string, v string) bool) {
}

// SpanID implements ddtrace.SpanContext.
func (*noopSpanContext) SpanID() uint64 {
	return 0
}

// TraceID implements ddtrace.SpanContext.
func (*noopSpanContext) TraceID() uint64 {
	return 0
}
