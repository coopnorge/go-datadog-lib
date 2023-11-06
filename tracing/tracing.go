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
// Deprecated: Use CreateChildSpan instead.
func CreateNestedTrace(sourceCtx context.Context, operation, resource string) (ddtrace.Span, error) {
	return CreateChildSpan(sourceCtx, operation, resource), nil
}

// CreateChildSpan will create a child-span of the span embedded in sourceCtx.
// If there is no trace-information in sourceCtx, a noop-span will be returned.
// The caller is responsible for calling span.Finish().
func CreateChildSpan(sourceCtx context.Context, operation, resource string) ddtrace.Span {
	existingSpan := getSpanFromContext(sourceCtx)
	if existingSpan == nil {
		return noopSpan{}
	}

	span := tracer.StartSpan(operation, tracer.ResourceName(resource), tracer.ChildOf(existingSpan.Context()))

	return span
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

// ExecuteWithTrace wraps a callable function with a Datadog tracer span.
// If the context.Context does not contain Datadog metadata, then
// the tracer span will not be recorded and the callable function
// will be executed as normal. Any error returned from the callable function
// will be recorded in the span, and then propagated to the caller of ExecuteWithTrace.
//
// Example:
//
//	err := ExecuteWithTrace(
//	    ctx,
//	    myAmazingFunctionToCall,
//	    "GetMyModel.SQL",
//	    "GetMyModel",
//	)
func ExecuteWithTrace(ctx context.Context, callable func(context.Context) error, source, op string) error {
	traceSpan := CreateChildSpan(ctx, fmt.Sprintf("%s.%s", source, op), op)
	traceCtx := addSpanToContext(ctx, traceSpan)

	execErr := callable(traceCtx)

	traceSpan.Finish(tracer.WithError(execErr))

	return execErr
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

func addSpanToContext(ctx context.Context, span tracer.Span) context.Context {
	if internal.IsExperimentalTracingEnabled() {
		return tracer.ContextWithSpan(ctx, span)
	}
	return internal.ExtendedContextWithMetadata(ctx, internal.TraceContextKey{}, TraceDetails{DatadogSpan: span})
}
