package tracing_test

import (
	"context"
	"testing"

	"github.com/coopnorge/go-datadog-lib/v2/tracing"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/mocktracer"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

func TestCreateNestedTrace(t *testing.T) {
	op := "test"
	res := "unit"
	ctx := context.Background()

	nestedTrace, nestedTraceErr := tracing.CreateNestedTrace(ctx, op, res)

	assert.NoError(t, nestedTraceErr)
	assert.NotNil(t, nestedTrace)
	// assert.IsType(t, noopSpan{}, nestedTrace)

	span, spanCtx := tracer.StartSpanFromContext(ctx, "test", tracer.ResourceName("UnitTest"))
	defer span.Finish()
	nestedTrace, nestedTraceErr = tracing.CreateNestedTrace(spanCtx, op, res)

	assert.Nil(t, nestedTraceErr)
	assert.NotNil(t, nestedTrace)
}

func TestAppendUserToTrace(t *testing.T) {
	// This test ensures that the legacy (deprecated) "AppendUserToTrace" no longer adds any personally identifiable information (PII) to the trace.

	// Start Datadog tracer, so that we don't create NoopSpans.
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

func TestResourceNameInTag(t *testing.T) {
	// This test ensures that the resource name is correctly set in the span tags.
	// Start Datadog tracer, so that we don't create NoopSpans.
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

func TestOverrideTraceResourceName(t *testing.T) {
	// Start Datadog tracer, so that we don't create NoopSpans.
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

func TestStartChildSpan(t *testing.T) {
	// Start Datadog tracer, so that we don't create NoopSpans.
	_ = mocktracer.Start()

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
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			if tt.args.spanInCtx {
				span, spanCtx := tracer.StartSpanFromContext(ctx, "test", tracer.ResourceName("UnitTest"))
				defer span.Finish()
				ctx = spanCtx
			}

			childSpan := tracing.CreateChildSpan(ctx, "my-operation", "my-resource")

			require.NotNil(t, childSpan)
			childSpan.Finish()
			if tt.args.spanInCtx {
				assert.NotEqual(t, uint64(0), childSpan.Context().SpanID())
				assert.NotEqual(t, uint64(0), childSpan.Context().TraceID())
			} else {
				assert.Equal(t, uint64(0), childSpan.Context().SpanID())
				assert.Equal(t, uint64(0), childSpan.Context().TraceID())
			}
		})
	}
}
