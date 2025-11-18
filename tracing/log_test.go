package tracing_test

import (
	"context"
	"testing"

	"github.com/DataDog/dd-trace-go/v2/ddtrace/tracer"
	"github.com/coopnorge/go-datadog-lib/v2/tracing"
	"github.com/coopnorge/go-logger"
)

func TestLogWithTrace(_ *testing.T) {
	ctx := context.Background()

	tracing.LogWithTrace(ctx, logger.LevelDebug, "unit test")
}

func TestLogFieldsWithTrace(_ *testing.T) {
	ctx := context.Background()

	tracing.LogFieldsWithTrace(ctx, logger.LevelDebug, "unit test", logger.Fields{})
}

func TestLogWithExtendedDatadogContext(_ *testing.T) {
	ctx := context.Background()
	span, spanCtx := tracer.StartSpanFromContext(ctx, "test", tracer.ResourceName("UnitTest"))
	defer span.Finish()
	tracing.LogWithTrace(spanCtx, logger.LevelDebug, "unit test")
}

func TestLogWithAllSeverity(_ *testing.T) {
	ctx := context.Background()

	tracing.LogWithTrace(ctx, logger.LevelDebug, "unit test")
	tracing.LogWithTrace(ctx, logger.LevelInfo, "unit test")
	tracing.LogWithTrace(ctx, logger.LevelWarn, "unit test")
	tracing.LogWithTrace(ctx, logger.LevelError, "unit test")
	// logger.LevelFatal will fail the test
	//LogWithTrace(ctx, logger.LevelFatal, "unit test")
}
