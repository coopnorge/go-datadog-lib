package coopdatadog

import (
	"context"
	"errors"
	"fmt"

	"github.com/coopnorge/go-datadog-lib/v2/internal"
	datadogLogger "github.com/coopnorge/go-logger/adapter/datadog"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
	"gopkg.in/DataDog/dd-trace-go.v1/profiler"
)

// Start the Datadog integration, use the returned Cancel function to stop the
// Datadog integration. When calling the StopFunc function traces will be
// flushed and profiling will be stopped to Datadog.
func Start(ctx context.Context, opts ...Option) (StopFunc, error) {
	if ctx == nil {
		return noop, errors.New("ctx cannot be nil")
	}
	if internal.IsDatadogDisabled() {
		return noop, nil
	}

	err := internal.VerifyEnvVarsSet(
		internal.DatadogAPMEndpoint,
		internal.DatadogDSDEndpoint,
		internal.DatadogService,
		internal.DatadogEnvironment,
		internal.DatadogVersion,
	)
	if err != nil {
		return noop, err
	}

	options, err := resolveOptions(opts)
	if err != nil {
		return noop, err
	}

	l, err := datadogLogger.NewLogger(datadogLogger.WithGlobalLogger())
	if err != nil {
		return noop, fmt.Errorf("Failed to initialize the Datadog logger: %w", err)
	}
	ddtrace.UseLogger(l)

	cancel := func() error {
		return stop(options)
	}

	err = start(options)
	return cancel, err
}

// StopFunc is a function signature for functions that stops the Datadog
// integration.
type StopFunc func() error

func noop() error {
	return nil
}

var _ StopFunc = noop

func start(options *options) error {
	startTracer()
	err := startProfiler(options)
	if err != nil {
		return err
	}
	return nil
}

func startTracer() {
	tracer.Start(
		tracer.WithRuntimeMetrics(),
	)
}

func startProfiler(options *options) error {
	var profilerTypes []profiler.ProfileType
	if options.enableExtraProfiling {
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
func stop(options *options) error {
	ctx := context.Background()

	if options.stopTimeout != nil {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, *options.stopTimeout)
		defer cancel()
	}

	errCh := make(chan error, 1)
	go func() {
		tracer.Stop()
		profiler.Stop()
		errCh <- nil
	}()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case err := <-errCh:
		return err
	}
}
