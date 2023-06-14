package tracing

import (
	"context"
	"testing"

	"github.com/coopnorge/go-datadog-lib/internal"

	"github.com/stretchr/testify/assert"
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
