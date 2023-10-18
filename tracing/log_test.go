package tracing

import (
	"context"
	"testing"

	"github.com/coopnorge/go-logger"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

func TestLogWithTrace(_ *testing.T) {
	ctx := context.Background()

	LogWithTrace(ctx, logger.LevelDebug, "unit test")
}

func TestLogFieldsWithTrace(_ *testing.T) {
	ctx := context.Background()

	LogFieldsWithTrace(ctx, logger.LevelDebug, "unit test", logger.Fields{})
}

func TestLogWithExtendedDatadogContext(_ *testing.T) {
	ctx := context.Background()
	span, spanCtx := tracer.StartSpanFromContext(ctx, "test", tracer.ResourceName("UnitTest"))
	defer span.Finish()
	LogWithTrace(spanCtx, logger.LevelDebug, "unit test")
}

func TestLogWithAllSeverity(_ *testing.T) {
	ctx := context.Background()

	LogWithTrace(ctx, logger.LevelDebug, "unit test")
	LogWithTrace(ctx, logger.LevelInfo, "unit test")
	LogWithTrace(ctx, logger.LevelWarn, "unit test")
	LogWithTrace(ctx, logger.LevelError, "unit test")
	// logger.LevelFatal will fail the test
	//LogWithTrace(ctx, logger.LevelFatal, "unit test")
}
