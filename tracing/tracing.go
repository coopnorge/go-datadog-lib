package tracing

import (
	"context"
	"fmt"

	"github.com/coopnorge/go-datadog-lib/internal"

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
	ddCtx, ddExist := internal.GetContextMetadata[TraceDetails](sourceCtx, internal.TraceContextKey{})
	if !ddExist || ddCtx.DatadogSpan == nil {
		return nil, fmt.Errorf("inheritance failed, parent span tracer not found in context")
	}

	nestedSpan := tracer.StartSpan(operation, tracer.ResourceName(resource), tracer.ChildOf(ddCtx.DatadogSpan.Context()))

	return nestedSpan, nil
}

// AppendUserToTrace includes identifier of user that would be attached to span in datadog
func AppendUserToTrace(sourceCtx context.Context, user string) error {
	ddCtx, ddExist := internal.GetContextMetadata[TraceDetails](sourceCtx, internal.TraceContextKey{})
	if !ddExist || ddCtx.DatadogSpan == nil {
		return fmt.Errorf("parent span tracer not found in context")
	}

	tracer.SetUser(ddCtx.DatadogSpan, user)

	return nil
}

// OverrideTraceResourceName set custom resource name for traced span aka SQL Query, Request, I/O etc
func OverrideTraceResourceName(sourceCtx context.Context, newResourceName string) error {
	ddCtx, ddExist := internal.GetContextMetadata[TraceDetails](sourceCtx, internal.TraceContextKey{})
	if !ddExist || ddCtx.DatadogSpan == nil {
		return fmt.Errorf("parent span tracer not found in context")
	}

	ddCtx.DatadogSpan.SetTag(ext.ResourceName, newResourceName)

	return nil
}
