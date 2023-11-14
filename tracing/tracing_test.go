package tracing

import (
	"context"
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
	t.Setenv(internal.ExperimentalTracingEnabled, "true")
	op := "test"
	res := "unit"
	ctx := context.Background()

	nestedTrace, nestedTraceErr := CreateNestedTrace(ctx, op, res)

	assert.Error(t, nestedTraceErr, "expected error since context not extended")
	assert.Nil(t, nestedTrace)

	span, spanCtx := tracer.StartSpanFromContext(ctx, "test", tracer.ResourceName("UnitTest"))
	defer span.Finish()

	nestedTrace, nestedTraceErr = CreateNestedTrace(spanCtx, op, res)

	assert.Nil(t, nestedTraceErr)
	assert.NotNil(t, nestedTrace)
}

func TestAppendUserToTrace(t *testing.T) {
	t.Setenv(internal.ExperimentalTracingEnabled, "false")
	// Start Datadog tracer, so that we don't create NoopSpans.
	testTracer := mocktracer.Start()
	t.Cleanup(testTracer.Stop)
	user := "unit_tester"
	ctx := context.Background()

	err := AppendUserToTrace(ctx, user)
	require.NoError(t, err)

	span, spanCtx := tracer.StartSpanFromContext(ctx, "test", tracer.ResourceName("UnitTest"))
	extCtx := internal.ExtendedContextWithMetadata(spanCtx, internal.TraceContextKey{}, TraceDetails{DatadogSpan: span})
	err = AppendUserToTrace(extCtx, user)
	require.NoError(t, err)
	span.Finish()

	testTracer.Stop()

	spans := testTracer.FinishedSpans()
	require.Equal(t, 1, len(spans))
	finishedSpan := spans[0]
	tags := finishedSpan.Tags()
	require.Equal(t, 1, len(tags), tags)
	require.Equal(t, "UnitTest", tags["resource.name"])
	require.Empty(t, tags["usr"])
	require.Empty(t, tags["usr.id"])
}

func TestAppendUserToTraceExperimental(t *testing.T) {
	t.Setenv(internal.ExperimentalTracingEnabled, "true")
	// Start Datadog tracer, so that we don't create NoopSpans.
	testTracer := mocktracer.Start()
	t.Cleanup(testTracer.Stop)
	user := "unit_tester"
	ctx := context.Background()

	err := AppendUserToTrace(ctx, user)

	require.NoError(t, err)

	span, spanCtx := tracer.StartSpanFromContext(ctx, "test", tracer.ResourceName("UnitTest"))
	err = AppendUserToTrace(spanCtx, user)
	require.NoError(t, err)
	span.Finish()

	testTracer.Stop()

	spans := testTracer.FinishedSpans()
	require.Equal(t, 1, len(spans))
	finishedSpan := spans[0]
	tags := finishedSpan.Tags()
	require.Equal(t, 1, len(tags), tags)
	require.Equal(t, "UnitTest", tags["resource.name"])
	require.Empty(t, tags["usr"])
	require.Empty(t, tags["usr.id"])
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

	require.NoError(t, err)
}

func TestOverrideTraceResourceNameExperimental(t *testing.T) {
	t.Setenv(internal.ExperimentalTracingEnabled, "true")
	newRes := "unit_test"
	ctx := context.Background()

	err := OverrideTraceResourceName(ctx, newRes)

	assert.Error(t, err, "expected error since context not extended")

	span, spanCtx := tracer.StartSpanFromContext(ctx, "test", tracer.ResourceName("UnitTest"))
	defer span.Finish()
	err = OverrideTraceResourceName(spanCtx, newRes)

	require.NoError(t, err)
}
