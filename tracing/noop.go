package tracing

import "gopkg.in/DataDog/dd-trace-go.v1/ddtrace"

// noopSpan and noopSpanContext is duplicated from "gopkg.in/DataDog/dd-trace-go.v1/ddtrace/internal/globaltracer.go"

var _ ddtrace.Span = (*noopSpan)(nil)

// noopSpan is an implementation of ddtrace.Span that is a no-op.
type noopSpan struct{}

// SetTag implements ddtrace.Span.
func (noopSpan) SetTag(_ string, _ interface{}) {}

// SetOperationName implements ddtrace.Span.
func (noopSpan) SetOperationName(_ string) {}

// BaggageItem implements ddtrace.Span.
func (noopSpan) BaggageItem(_ string) string { return "" }

// SetBaggageItem implements ddtrace.Span.
func (noopSpan) SetBaggageItem(_, _ string) {}

// Finish implements ddtrace.Span.
func (noopSpan) Finish(_ ...ddtrace.FinishOption) {}

// Context implements ddtrace.Span.
func (noopSpan) Context() ddtrace.SpanContext { return noopSpanContext{} }

var _ ddtrace.SpanContext = (*noopSpanContext)(nil)

// noopSpanContext is an implementation of ddtrace.SpanContext that is a no-op.
type noopSpanContext struct{}

// SpanID implements ddtrace.SpanContext.
func (noopSpanContext) SpanID() uint64 { return 0 }

// TraceID implements ddtrace.SpanContext.
func (noopSpanContext) TraceID() uint64 { return 0 }

// ForeachBaggageItem implements ddtrace.SpanContext.
func (noopSpanContext) ForeachBaggageItem(_ func(k, v string) bool) {}
