package tracing

import (
	"context"

	"github.com/coopnorge/go-datadog-lib/v2/internal"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/ext"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

type (
	// TraceDetails that must be included in context
	TraceDetails struct {
		// DatadogSpan represents a chunk of computation time for Datadog system
		DatadogSpan ddtrace.Span
	}
)

// CreateNestedTrace will fork parent tracer to attach to parent one with new operation and resource from sourceCtx
func CreateNestedTrace(sourceCtx context.Context, operation, resource string) (ddtrace.Span, error) {
	span := GetSpanFromContext(sourceCtx)

	nestedSpan := tracer.StartSpan(operation, tracer.ResourceName(resource), tracer.ChildOf(span.Context()))

	return nestedSpan, nil
}

// AppendUserToTrace includes identifier of user that would be attached to span in datadog
func AppendUserToTrace(sourceCtx context.Context, user string) error {
	span := GetSpanFromContext(sourceCtx)

	tracer.SetUser(span, user)

	return nil
}

// OverrideTraceResourceName set custom resource name for traced span aka SQL Query, Request, I/O etc
func OverrideTraceResourceName(sourceCtx context.Context, newResourceName string) error {
	span := GetSpanFromContext(sourceCtx)

	span.SetTag(ext.ResourceName, newResourceName)

	return nil
}

func GetSpanFromContext(ctx context.Context) tracer.Span {
	if internal.IsExperimentalTracingEnabled() {
		if span, exists := tracer.SpanFromContext(ctx); exists {
			return span
		}
		return &noopSpan{}
	}
	ddCtx, ddExist := internal.GetContextMetadata[TraceDetails](ctx, internal.TraceContextKey{})
	if !ddExist || ddCtx.DatadogSpan == nil {
		return &noopSpan{}
	}
	return ddCtx.DatadogSpan
}
