package tracing

import (
	"context"
	"os"
	"sort"
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

	assert.Error(t, nestedTraceErr, "expected error since context not extended")
	assert.Nil(t, nestedTrace)

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

	assert.Error(t, nestedTraceErr, "expected error since context not extended")
	assert.Nil(t, nestedTrace)

	span, spanCtx := tracer.StartSpanFromContext(ctx, "test", tracer.ResourceName("UnitTest"))
	defer span.Finish()
	os.Setenv(internal.ExperimentalTracingEnabled, "true")
	defer func() {
		os.Setenv(internal.ExperimentalTracingEnabled, "")
	}()
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
	os.Setenv(internal.ExperimentalTracingEnabled, "true")
	defer func() {
		os.Setenv(internal.ExperimentalTracingEnabled, "")
	}()
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
	os.Setenv(internal.ExperimentalTracingEnabled, "true")
	defer func() {
		os.Setenv(internal.ExperimentalTracingEnabled, "")
	}()
	err = OverrideTraceResourceName(spanCtx, newRes)

	assert.Nil(t, err)
}

func TestExecuteWithTrace(t *testing.T) {
	ctx := context.Background()

	// Start Datadog tracer, so that we don't create NoopSpans.
	mocktracer := mocktracer.Start()

	span, spanCtx := tracer.StartSpanFromContext(ctx, "test", tracer.ResourceName("UnitTest"))
	defer span.Finish()

	ddCtx := internal.ExtendedContextWithMetadata(spanCtx, internal.TraceContextKey{}, TraceDetails{DatadogSpan: span})

	isCallableCalled := false
	callableUnit := func(ctx context.Context) (bool, error) {
		metadata, ok := internal.GetContextMetadata[TraceDetails](ctx, internal.TraceContextKey{})
		require.True(t, ok)
		assert.Equal(t, span.Context().TraceID(), metadata.DatadogSpan.Context().TraceID(), "The TraceID of both spans should be the same (same root)")
		assert.NotEqual(t, span.Context().SpanID(), metadata.DatadogSpan.Context().SpanID(), "The SpanID should be different, and it should be a child-span")
		return !isCallableCalled, nil
	}

	execRes, execErr := ExecuteWithTrace[bool](ddCtx, callableUnit, "unit.test", "test")
	assert.NoError(t, execErr, "expected context with Datadog tracer context")
	assert.True(t, execRes, "expected response of ExecuteWithTrace to be true")

	span.Finish() // Finish the original span

	require.Equal(t, 0, len(mocktracer.OpenSpans()))
	spans := mocktracer.FinishedSpans()
	sort.Slice(spans, func(i, j int) bool {
		return spans[i].StartTime().Before(spans[j].StartTime())
	})
	require.Equal(t, 2, len(spans))

	assert.NotEqual(t, uint64(0), spans[0].SpanID())
	assert.NotEqual(t, uint64(0), spans[1].SpanID())
	assert.Equal(t, uint64(0), spans[0].ParentID())
	assert.Equal(t, spans[0].SpanID(), spans[1].ParentID(), "spans[1] should be child of spans[0]")
}
