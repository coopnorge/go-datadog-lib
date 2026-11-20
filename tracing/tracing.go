// Package tracing creates and manipulates Datadog spans.
//
// Deprecated: The whole module is deprecated. Use otel.Tracer(name).Start from
// go.opentelemetry.io/otel for spans. The loghook in
// github.com/coopnorge/go-telemetry-lib/loghook adds trace and span IDs to log
// entries. See the root package documentation for the migration guide and ADR.
package tracing

import (
	"context"
	"fmt"

	"github.com/DataDog/dd-trace-go/v2/ddtrace/ext"
	"github.com/DataDog/dd-trace-go/v2/ddtrace/tracer"
)

type (
	// TraceDetails that must be included in context
	//
	// Deprecated: The whole module is deprecated. Use otel.Tracer(name).Start from
	// go.opentelemetry.io/otel; nesting happens through the context.
	TraceDetails struct {
		// DatadogSpan represents a chunk of computation time for Datadog system
		//
		// Deprecated: The whole module is deprecated. Use otel.Tracer(name).Start from
		// go.opentelemetry.io/otel; nesting happens through the context.
		DatadogSpan *tracer.Span
	}
)

// CreateNestedTrace will fork parent tracer to attach to parent one with new operation and resource from the provided context.
//
// Deprecated: Use CreateChildSpan instead.
// The whole module is deprecated; use otel.Tracer(name).Start from
// go.opentelemetry.io/otel; nesting happens through the context.
func CreateNestedTrace(ctx context.Context, operation, resource string) (*tracer.Span, error) {
	return CreateChildSpan(ctx, operation, resource), nil
}

// CreateChildSpan will create a child-span of the span embedded in the provided context.
// If there is no trace-information in the provided context, a noop-span (nil) will be returned.
// The caller is responsible for calling span.Finish().
//
// Deprecated: The whole module is deprecated. Use otel.Tracer(name).Start from
// go.opentelemetry.io/otel; nesting happens through the context.
func CreateChildSpan(ctx context.Context, operation, resource string) *tracer.Span {
	existingSpan, _ := tracer.SpanFromContext(ctx)
	return existingSpan.StartChild(operation, tracer.ResourceName(resource))
}

// AppendUserToTrace includes identifier of user that would be attached to span in datadog
//
// Deprecated: AppendUserToTrace previously added CoopID to Datadog-spans, which could be used to look up other PII-information. This is not wanted, and has been replaced with a no-op.
// The whole module is deprecated. There is no direct replacement.
func AppendUserToTrace(_ context.Context, _ string) error {
	return nil
}

// OverrideTraceResourceName set custom resource name for traced span aka SQL Query, Request, I/O etc
//
// Deprecated: The whole module is deprecated. There is no direct replacement.
func OverrideTraceResourceName(ctx context.Context, newResourceName string) error {
	span, exists := tracer.SpanFromContext(ctx)
	if !exists {
		return fmt.Errorf("parent span tracer not found in context")
	}

	span.SetTag(ext.ResourceName, newResourceName)

	return nil
}
