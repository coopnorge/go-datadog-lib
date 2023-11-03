package tracing

import (
	"context"
	"fmt"
	"testing"

	"github.com/coopnorge/go-datadog-lib/v2/internal"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/mocktracer"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

func TestCreateNestedTrace(t *testing.T) {
	op := "test"
	res := "unit"
	ctx := context.Background()

	nestedTrace, nestedTraceErr := CreateNestedTrace(ctx, op, res)

	assert.NoError(t, nestedTraceErr)
	assert.IsType(t, nestedTrace, noopSpan{})

	span, spanCtx := tracer.StartSpanFromContext(ctx, "test", tracer.ResourceName("UnitTest"))
	defer span.Finish()
	extCtx := internal.ExtendedContextWithMetadata(spanCtx, internal.TraceContextKey{}, TraceDetails{DatadogSpan: span})
	nestedTrace, nestedTraceErr = CreateNestedTrace(extCtx, op, res)

	assert.Nil(t, nestedTraceErr)
	assert.NotNil(t, nestedTrace)
}

func TestCreateNestedTraceExperimental(t *testing.T) {
	op := "test"
	res := "unit"
	ctx := context.Background()

	nestedTrace, nestedTraceErr := CreateNestedTrace(ctx, op, res)

	assert.NoError(t, nestedTraceErr)
	assert.IsType(t, nestedTrace, noopSpan{})

	span, spanCtx := tracer.StartSpanFromContext(ctx, "test", tracer.ResourceName("UnitTest"))
	defer span.Finish()
	t.Setenv(internal.ExperimentalTracingEnabled, "true")
	nestedTrace, nestedTraceErr = CreateNestedTrace(spanCtx, op, res)

	assert.Nil(t, nestedTraceErr)
	assert.NotNil(t, nestedTrace)
}

func TestAppendUserToTrace(t *testing.T) {
	user := "unit_tester"
	ctx := context.Background()

	err := AppendUserToTrace(ctx, user)

	assert.Error(t, err, "expected error since context not extended")

	span, spanCtx := tracer.StartSpanFromContext(ctx, "test", tracer.ResourceName("UnitTest"))
	defer span.Finish()
	extCtx := internal.ExtendedContextWithMetadata(spanCtx, internal.TraceContextKey{}, TraceDetails{DatadogSpan: span})
	err = AppendUserToTrace(extCtx, user)

	assert.Nil(t, err)
}

func TestAppendUserToTraceExperimental(t *testing.T) {
	user := "unit_tester"
	ctx := context.Background()

	err := AppendUserToTrace(ctx, user)

	assert.Error(t, err, "expected error since context not extended")

	span, spanCtx := tracer.StartSpanFromContext(ctx, "test", tracer.ResourceName("UnitTest"))
	defer span.Finish()
	t.Setenv(internal.ExperimentalTracingEnabled, "true")
	err = AppendUserToTrace(spanCtx, user)

	assert.Nil(t, err)
}

func TestOverrideTraceResourceName(t *testing.T) {
	newRes := "unit_test"
	ctx := context.Background()

	err := OverrideTraceResourceName(ctx, newRes)

	assert.Error(t, err, "expected error since context not extended")

	span, spanCtx := tracer.StartSpanFromContext(ctx, "test", tracer.ResourceName("UnitTest"))
	defer span.Finish()
	extCtx := internal.ExtendedContextWithMetadata(spanCtx, internal.TraceContextKey{}, TraceDetails{DatadogSpan: span})
	err = OverrideTraceResourceName(extCtx, newRes)

	assert.Nil(t, err)
}

func TestOverrideTraceResourceNameExperimental(t *testing.T) {
	newRes := "unit_test"
	ctx := context.Background()

	err := OverrideTraceResourceName(ctx, newRes)

	assert.Error(t, err, "expected error since context not extended")

	span, spanCtx := tracer.StartSpanFromContext(ctx, "test", tracer.ResourceName("UnitTest"))
	defer span.Finish()
	t.Setenv(internal.ExperimentalTracingEnabled, "true")
	err = OverrideTraceResourceName(spanCtx, newRes)

	assert.Nil(t, err)
}

func TestStartChildSpan(t *testing.T) {
	// Start Datadog tracer, so that we don't create NoopSpans.
	_ = mocktracer.Start()

	type args struct {
		spanInCtx    bool
		experimental bool
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "legacy without span",
			args: args{
				spanInCtx:    false,
				experimental: false,
			},
		},
		{
			name: "legacy with span",
			args: args{
				spanInCtx:    true,
				experimental: false,
			},
		},
		{
			name: "experimental without span",
			args: args{
				spanInCtx:    false,
				experimental: true,
			},
		},
		{
			name: "experimental with span",
			args: args{
				spanInCtx:    true,
				experimental: true,
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Setenv(internal.ExperimentalTracingEnabled, fmt.Sprintf("%v", tt.args.experimental))

			ctx := context.Background()
			if tt.args.spanInCtx {
				if tt.args.experimental {
					span, spanCtx := tracer.StartSpanFromContext(ctx, "test", tracer.ResourceName("UnitTest"))
					defer span.Finish()
					ctx = spanCtx
				} else {
					span, spanCtx := tracer.StartSpanFromContext(ctx, "test", tracer.ResourceName("UnitTest"))
					defer span.Finish()
					extCtx := internal.ExtendedContextWithMetadata(spanCtx, internal.TraceContextKey{}, TraceDetails{DatadogSpan: span})
					ctx = extCtx
				}
			}

			childSpan := CreateChildSpan(ctx, "my-operation", "my-resource")

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
