package tracing

import "gopkg.in/DataDog/dd-trace-go.v1/ddtrace"

type (
	noopSpan struct{}
)

var _ ddtrace.Span = (*noopSpan)(nil)

// BaggageItem implements ddtrace.Span.
func (*noopSpan) BaggageItem(_ string) string {
	return ""
}

// Context implements ddtrace.Span.
func (*noopSpan) Context() ddtrace.SpanContext {
	return &noopSpanContext{}
}

// Finish implements ddtrace.Span.
func (*noopSpan) Finish(_ ...ddtrace.FinishOption) {
}

// SetBaggageItem implements ddtrace.Span.
func (*noopSpan) SetBaggageItem(_ string, _ string) {
}

// SetOperationName implements ddtrace.Span.
func (*noopSpan) SetOperationName(_ string) {
}

// SetTag implements ddtrace.Span.
func (*noopSpan) SetTag(_ string, _ interface{}) {
}
