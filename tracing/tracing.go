package tracing

import (
	"context"
	"fmt"

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
	span := getSpanFromContext(sourceCtx)
	if span == nil {
		return nil, fmt.Errorf("inheritance failed, parent span tracer not found in context")
	}

	nestedSpan := tracer.StartSpan(operation, tracer.ResourceName(resource), tracer.ChildOf(span.Context()))

	return nestedSpan, nil
}

// AppendUserToTrace includes identifier of user that would be attached to span in datadog
func AppendUserToTrace(sourceCtx context.Context, user string) error {
	span := getSpanFromContext(sourceCtx)
	if span == nil {
		return fmt.Errorf("parent span tracer not found in context")
	}

	tracer.SetUser(span, user)

	return nil
}

// OverrideTraceResourceName set custom resource name for traced span aka SQL Query, Request, I/O etc
func OverrideTraceResourceName(sourceCtx context.Context, newResourceName string) error {
	span := getSpanFromContext(sourceCtx)
	if span == nil {
		return fmt.Errorf("parent span tracer not found in context")
	}

	span.SetTag(ext.ResourceName, newResourceName)

	return nil
}

func getSpanFromContext(ctx context.Context) tracer.Span {
	if internal.IsExperimentalTracingEnabled() {
		if span, exists := tracer.SpanFromContext(ctx); exists {
			return span
		}
		return nil
	}
	ddCtx, ddExist := internal.GetContextMetadata[TraceDetails](ctx, internal.TraceContextKey{})
	if !ddExist || ddCtx.DatadogSpan == nil {
		return nil
	}
	return ddCtx.DatadogSpan
}
