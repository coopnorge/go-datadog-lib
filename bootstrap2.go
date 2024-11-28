package coopdatadog

import (
	"context"
	"fmt"

	"github.com/coopnorge/go-datadog-lib/v2/internal"
	datadogLogger "github.com/coopnorge/go-logger/adapter/datadog"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
	"gopkg.in/DataDog/dd-trace-go.v1/profiler"
)

// Start the Datadog integration, use the returned context.CancelFunc to stop
// the Datadog integration.
func Start(ctx context.Context, options ...Option) (context.Context, context.CancelFunc, error) {
	if internal.IsDatadogDisabled() {
		return ctx, context.CancelFunc(func() {}), nil
	}

	err := internal.VerifyEnvVarsSet(
		internal.DatadogAPMEndpoint,
		internal.DatadogDSDEndpoint,
		internal.DatadogService,
		internal.DatadogEnvironment,
		internal.DatadogVersion,
	)
	if err != nil {
		return ctx, context.CancelFunc(func() {}), err
	}

	cfg, err := resolveConfig(options)
	if err != nil {
		return ctx, context.CancelFunc(func() {}), err
	}

	ctx, cancel := context.WithCancel(ctx)
	canceldd := context.CancelFunc(func() {
		stop(cfg)
		cancel()
	})

	l, err := datadogLogger.NewLogger(datadogLogger.WithGlobalLogger())
	if err != nil {
		return nil, nil, fmt.Errorf("Failed to initialize the Datadog logger: %w", err)
	}
	ddtrace.UseLogger(l)

	err = start(ctx, cfg)
	return ctx, canceldd, err
}

func start(_ context.Context, cfg *config) error {
	startTracer(cfg)
	err := startProfiler(cfg)
	if err != nil {
		return err
	}
	return nil
}

func startTracer(cfg *config) {
	if !cfg.enableTracing {
		return
	}
	tracer.Start(
		tracer.WithRuntimeMetrics(),
	)
}

func startProfiler(cfg *config) error {
	if !cfg.enableProfiling {
		return nil
	}

	var profilerTypes []profiler.ProfileType
	if cfg.enableExtraProfiling {
		profilerTypes = []profiler.ProfileType{
			profiler.CPUProfile,
			profiler.HeapProfile,
			profiler.GoroutineProfile,
			profiler.MutexProfile,
			profiler.BlockProfile,
		}
	} else {
		profilerTypes = []profiler.ProfileType{profiler.CPUProfile}
	}

	return profiler.Start(profiler.WithProfileTypes(profilerTypes...))
}

// stop with a graceful shutdown that includes flushing signals.
func stop(_ *config) {
	tracer.Stop()
	profiler.Stop()
}
