package coopdatadog

import (
	"fmt"

	"github.com/coopnorge/go-datadog-lib/v2/internal"
	datadogLogger "github.com/coopnorge/go-logger/adapter/datadog"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
	"gopkg.in/DataDog/dd-trace-go.v1/profiler"
)

// Start the Datadog integration, use the returned Cancel function to stop the
// Datadog integration. When calling the Cancel function traces will be flushed
// and profiling will be stopped to Datadog.
//
// Usage:
//
//	cancel, err := Start()
//	if err != nil {
//		panic(err)
//	}
//	defer func() {
//		err := cancel()
//		if err != nil {
//			panic(err)
//		}
//	}}
func Start(opts ...Option) (Cancel, error) {
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
		return stop()
	}

	err = start(options)
	return cancel, err
}

// Cancel is a function signature for functions that stops the Datadog
// integration.
type Cancel func() error

func noop() error {
	return nil
}

var _ Cancel = noop

func start(options *options) error {
	startTracer(options)
	err := startProfiler(options)
	if err != nil {
		return err
	}
	return nil
}

func startTracer(options *options) {
	if !options.enableTracing {
		return
	}
	tracer.Start(
		tracer.WithRuntimeMetrics(),
	)
}

func startProfiler(options *options) error {
	if !options.enableProfiling {
		return nil
	}

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
func stop() error {
	tracer.Stop()
	profiler.Stop()
	return nil
}
