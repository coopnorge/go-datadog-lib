package tracelogger_test

import (
	"context"

	"github.com/coopnorge/go-datadog-lib/v2/tracelogger"
	"github.com/coopnorge/go-logger"

	"github.com/DataDog/dd-trace-go/v2/ddtrace/tracer"
)

func ExampleHook() {
	logger.ConfigureGlobalLogger(
		logger.WithHook(tracelogger.NewHook()),
	)
	ctx := context.Background()
	a(ctx)
}

func a(ctx context.Context) {
	span, ctx := tracer.StartSpanFromContext(ctx, "a")
	err := b(ctx)
	span.Finish(tracer.WithError(err))
}

func b(ctx context.Context) error {
	logger.WithContext(ctx).Info("Hello")
	return nil
}
