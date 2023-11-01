package tracing

import "gopkg.in/DataDog/dd-trace-go.v1/ddtrace"

type (
	noopSpan struct{}
)

var _ ddtrace.Span = (*noopSpan)(nil)

// BaggageItem implements ddtrace.Span.
func (*noopSpan) BaggageItem(key string) string {
	return ""
}

// Context implements ddtrace.Span.
func (*noopSpan) Context() ddtrace.SpanContext {
	return &noopSpanContext{}
}

// Finish implements ddtrace.Span.
func (*noopSpan) Finish(opts ...ddtrace.FinishOption) {
}

// SetBaggageItem implements ddtrace.Span.
func (*noopSpan) SetBaggageItem(key string, val string) {
}

// SetOperationName implements ddtrace.Span.
func (*noopSpan) SetOperationName(operationName string) {
}

// SetTag implements ddtrace.Span.
func (*noopSpan) SetTag(key string, value interface{}) {
}
