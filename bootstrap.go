package coopdatadog

import (
    "fmt"
    
    "github.com/coopnorge/go-datadog-lib/config"

    "gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
    "gopkg.in/DataDog/dd-trace-go.v1/profiler"
)

// StartDatadog parallel process to collect data for Datadog.
// enableExtraProfiling flag enables more optional profilers not recommended for production.
// isConnectionSocket flag related to Datadog connection type, it supports HTTP or socket - values will be used from config.DatadogParameters
func StartDatadog(cfg config.DatadogParameters, enableExtraProfiling, isConnectionSocket bool) error {
    if !cfg.IsDataDogConfigValid() {
        return fmt.Errorf("datadog configuration not valid, cannot initialize Datadog services")
    }

    initTracer(cfg, isConnectionSocket)
    if initProfilerErr := initProfiler(cfg, enableExtraProfiling, isConnectionSocket); initProfilerErr != nil {
        return fmt.Errorf("failed to start Datadog profiler: %v", initProfilerErr)
    }

    return nil
}

// GracefulDatadogShutdown of executed parallel processes
func GracefulDatadogShutdown() {
    defer tracer.Stop()
    defer profiler.Stop()
}

func initTracer(cfg config.DatadogParameters, isConnectionSocket bool) {
    var tracerOptions []tracer.StartOption
    if isConnectionSocket {
        tracerOptions = append(tracerOptions, tracer.WithUDS(cfg.GetApmEndpoint()))
    } else {
        tracerOptions = append(tracerOptions, tracer.WithAgentAddr(cfg.GetApmEndpoint()))
    }

    tracerOptions = append(
        tracerOptions,
        []tracer.StartOption{
            tracer.WithEnv(cfg.GetEnv()),
            tracer.WithService(cfg.GetService()),
            tracer.WithServiceVersion(cfg.GetServiceVersion()),
        }...,
    )

    tracer.Start(tracerOptions...)
}

func initProfiler(cfg config.DatadogParameters, enableExtraProfiling bool, isConnectionSocket bool) error {
    var profilerTypes []profiler.ProfileType
    if enableExtraProfiling {
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
    if isConnectionSocket {
        profilerOptions = append(profilerOptions, profiler.WithUDS(cfg.GetApmEndpoint()))
    } else {
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
