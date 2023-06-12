package coopdatadog

import (
	"fmt"

	"github.com/coopnorge/go-datadog-lib/config"

	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
	"gopkg.in/DataDog/dd-trace-go.v1/profiler"
)

// ConnectionType enum type
type ConnectionType byte

const (
	// ConnectionTypeSocket sets the connection to Datadog to go throug a UNIX socket
	ConnectionTypeSocket ConnectionType = iota
	// ConnectionTypeHTTP sets the connection to Datadog to go over HTTP
	ConnectionTypeHTTP
)

// StartDatadog parallel process to collect data for Datadog.
// connectionType flag related to Datadog connection type, it supports HTTP or socket - values will be used from config.DatadogParameters
func StartDatadog(cfg config.DatadogParameters, connectionType ConnectionType) error {
	if !cfg.IsDataDogConfigValid() {
		return fmt.Errorf("datadog configuration not valid, cannot initialize Datadog services")
	}

	initTracer(cfg, connectionType)
	if initProfilerErr := initProfiler(cfg, connectionType); initProfilerErr != nil {
		return fmt.Errorf("failed to start Datadog profiler: %v", initProfilerErr)
	}

	return nil
}

// GracefulDatadogShutdown of executed parallel processes
func GracefulDatadogShutdown() {
	defer tracer.Stop()
	defer profiler.Stop()
}

func initTracer(cfg config.DatadogParameters, connectionType ConnectionType) {
	var tracerOptions []tracer.StartOption
	switch connectionType {
	case ConnectionTypeSocket:
		tracerOptions = append(tracerOptions, tracer.WithUDS(cfg.GetApmEndpoint()))
	case ConnectionTypeHTTP:
		tracerOptions = append(tracerOptions, tracer.WithAgentAddr(cfg.GetApmEndpoint()))
	}

	tracerOptions = append(
		tracerOptions,
		[]tracer.StartOption{
			tracer.WithEnv(cfg.GetEnv()),
			tracer.WithService(cfg.GetService()),
			tracer.WithServiceVersion(cfg.GetServiceVersion()),
			tracer.WithRuntimeMetrics(),
		}...,
	)

	tracer.Start(tracerOptions...)
}

func initProfiler(cfg config.DatadogParameters, connectionType ConnectionType) error {
	var profilerTypes []profiler.ProfileType
	if cfg.IsExtraProfilingEnabled() {
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

	var profilerOptions []profiler.Option
	switch connectionType {
	case ConnectionTypeSocket:
		profilerOptions = append(profilerOptions, profiler.WithUDS(cfg.GetApmEndpoint()))
	case ConnectionTypeHTTP:
		profilerOptions = append(profilerOptions, profiler.WithAgentAddr(cfg.GetApmEndpoint()))
	}

	profilerOptions = append(
		profilerOptions,
		[]profiler.Option{
			profiler.WithEnv(cfg.GetEnv()),
			profiler.WithService(cfg.GetService()),
			profiler.WithVersion(cfg.GetServiceVersion()),
			profiler.WithProfileTypes(profilerTypes...),
		}...,
	)

	return profiler.Start(profilerOptions...)
}
