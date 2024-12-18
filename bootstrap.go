package coopdatadog

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/coopnorge/go-datadog-lib/v2/internal"
	datadogLogger "github.com/coopnorge/go-logger/adapter/datadog"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
	"gopkg.in/DataDog/dd-trace-go.v1/profiler"
)

// Start the Datadog integration, use the returned Cancel function to stop the
// Datadog integration. When calling the StopFunc function traces will be
// flushed and profiling will be stopped to Datadog.
//
// Usage:
//
//	package main
//
//	import (
//		"github.com/coopnorge/go-datadog-lib/v2"
//	)
//
//	func main() {
//		err := run()
//		if err != nil {
//			panic(err)
//		}
//	}
//
//	func run() error {
//		stop, err := coopdatadog.Start(context.Background())
//		if err != nil {
//			return err
//		}
//		defer func() {
//			err := stop()
//			if err != nil {
//				panic(err)
//			}
//		}()
//		return nil
//	}
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

	if err := normalizeDatadogEnvVars(); err != nil {
		return noop, fmt.Errorf("failed to normalize Datadog environment variables: %w", err)
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

// normalizeDatadogEnvVars ensures that the environment variables that the Datadog library is on the expected format.
func normalizeDatadogEnvVars() error {
	apmEndpoint := os.Getenv(internal.DatadogAPMEndpoint)
	if normalizedAPMEndpoint, changed := normalizeAPMEndpoint(apmEndpoint); changed {
		err := os.Setenv(internal.DatadogAPMEndpoint, normalizedAPMEndpoint)
		if err != nil {
			return err
		}
	}
	return nil
}

func normalizeAPMEndpoint(apmEndpoint string) (string, bool) {
	if strings.HasPrefix(apmEndpoint, "/") {
		// apmEndpoint did not have a scheme set, but it looks like unix scheme, so we explicitly set it.
		return fmt.Sprintf("unix://%s", apmEndpoint), true
	}
	return apmEndpoint, false
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
func stop() error {
	tracer.Stop()
	profiler.Stop()
	return nil
}
