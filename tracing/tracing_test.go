package tracing_test

import (
	"context"
	"testing"

	"github.com/DataDog/dd-trace-go/v2/ddtrace/mocktracer"
	"github.com/DataDog/dd-trace-go/v2/ddtrace/tracer"
	"github.com/coopnorge/go-datadog-lib/v2/tracing"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestCreateNestedTrace tests the CreateNestedTrace function with and without an existing span in the context.
func TestCreateNestedTrace(t *testing.T) {
	op := "test"
	res := "unit"
	ctx := context.Background()

	testTracer := mocktracer.Start()
	t.Cleanup(testTracer.Stop)

	nestedTrace, nestedTraceErr := tracing.CreateNestedTrace(ctx, op, res)

	assert.NoError(t, nestedTraceErr)
	assert.Nil(t, nestedTrace) // since we are using a context without a span, we get a nil span

	span, spanCtx := tracer.StartSpanFromContext(ctx, "test", tracer.ResourceName("UnitTest"))
	defer span.Finish()
	nestedTrace, nestedTraceErr = tracing.CreateNestedTrace(spanCtx, op, res)

	assert.Nil(t, nestedTraceErr)
	assert.NotNil(t, nestedTrace)
}

// TestAppendUserToTrace ensures that the legacy (deprecated) "AppendUserToTrace" no longer adds any personally identifiable information (PII) to the trace.
func TestAppendUserToTrace(t *testing.T) {
	testTracer := mocktracer.Start()
	t.Cleanup(testTracer.Stop)
	user := "unit_tester"
	ctx := context.Background()

	err := tracing.AppendUserToTrace(ctx, user)
	require.NoError(t, err)

	span, spanCtx := tracer.StartSpanFromContext(ctx, "test")
	err = tracing.AppendUserToTrace(spanCtx, user)
	require.NoError(t, err)
	span.Finish()

	testTracer.Stop()

	spans := testTracer.FinishedSpans()
	require.Equal(t, 1, len(spans))
	finishedSpan := spans[0]
	tags := finishedSpan.Tags()
	require.Empty(t, tags["usr"])
	require.Empty(t, tags["usr.id"])
}

// TestResourceNameInTag ensures that the resource name is correctly set in the span tags.
func TestResourceNameInTag(t *testing.T) {
	testTracer := mocktracer.Start()
	t.Cleanup(testTracer.Stop)

	span, _ := tracer.StartSpanFromContext(context.Background(), "test", tracer.ResourceName("UnitTest"))
	span.Finish()

	testTracer.Stop()

	spans := testTracer.FinishedSpans()
	require.Equal(t, 1, len(spans))
	finishedSpan := spans[0]
	tags := finishedSpan.Tags()
	require.Equal(t, "UnitTest", tags["resource.name"])
}

// TestOverrideTraceResourceName tests the OverrideTraceResourceName function with and without an existing span in the context.
func TestOverrideTraceResourceName(t *testing.T) {
	testTracer := mocktracer.Start()
	t.Cleanup(testTracer.Stop)

	newRes := "unit_test"
	ctx := context.Background()

	err := tracing.OverrideTraceResourceName(ctx, newRes)

	assert.Error(t, err, "expected error since context not extended")

	span, spanCtx := tracer.StartSpanFromContext(ctx, "test", tracer.ResourceName("UnitTest"))
	defer span.Finish()
	err = tracing.OverrideTraceResourceName(spanCtx, newRes)

	assert.Nil(t, err)
}

// TestStartChildSpan tests the CreateChildSpan function with and without an existing span in the context.
func TestStartChildSpan(t *testing.T) {
	testTracer := mocktracer.Start()
	t.Cleanup(testTracer.Stop)

	type args struct {
		spanInCtx bool
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "without span",
			args: args{
				spanInCtx: false,
			},
		},
		{
			name: "with span",
			args: args{
				spanInCtx: true,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			if tt.args.spanInCtx {
				span, spanCtx := tracer.StartSpanFromContext(ctx, "test", tracer.ResourceName("UnitTest"))
				defer span.Finish()
				ctx = spanCtx
			}

			childSpan := tracing.CreateChildSpan(ctx, "my-operation", "my-resource")

			if tt.args.spanInCtx {
				require.NotNil(t, childSpan)
			} else {
				require.Nil(t, childSpan) // since we are using a context without a span, we get a nil span
			}
			childSpan.Finish()
			if tt.args.spanInCtx {
				assert.NotEqual(t, uint64(0), childSpan.Context().SpanID())
				assert.NotEqual(t, uint64(0), childSpan.Context().TraceIDLower())
			} else {
				assert.Equal(t, uint64(0), childSpan.Context().SpanID())
				assert.Equal(t, uint64(0), childSpan.Context().TraceIDLower())
			}
		})
	}
}
