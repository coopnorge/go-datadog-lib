package tracing

import (
	"context"
	"testing"

	"github.com/coopnorge/go-logger"
)

func TestLogWithTrace(t *testing.T) {
	ctx := context.Background()

	LogWithTrace(ctx, logger.LevelDebug, "unit test")
}

func TestLogFieldsWithTrace(t *testing.T) {
	ctx := context.Background()

	LogFieldsWithTrace(ctx, logger.LevelDebug, "unit test", logger.Fields{})
}

func TestLogWithAllSeverity(t *testing.T) {
	ctx := context.Background()

	LogWithTrace(ctx, logger.LevelDebug, "unit test")
	LogWithTrace(ctx, logger.LevelInfo, "unit test")
	LogWithTrace(ctx, logger.LevelWarn, "unit test")
	LogWithTrace(ctx, logger.LevelError, "unit test")
	// logger.LevelFatal will fail the test
}
