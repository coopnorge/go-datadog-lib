package go_datadog_lib

import (
	"github.com/coopnorge/go-datadog-lib/config"
	"github.com/coopnorge/go-logger"

	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
	"gopkg.in/DataDog/dd-trace-go.v1/profiler"
)

// StartDatadog parallel process to collect data, enableExtraProfiling flag enables more optional profilers not recommended for prod
func StartDatadog(cfg config.DatadogConfig, enableExtraProfiling bool) {
	if config.IsDataDogConfigValid(cfg) {
		logger.Errorf("Datadog configuration not valid, cannot initialize Datadog services")

		return
	}

	logger.Infof("Initializing Datadog services for %s in environment %s", cfg.Service, cfg.Env)

	tracer.Start(
		tracer.WithEnv(cfg.Env),
		tracer.WithUDS(cfg.APM),
		tracer.WithService(cfg.Service),
		tracer.WithServiceVersion(cfg.ServiceVersion),
	)

	var profileTypes []profiler.ProfileType
	if enableExtraProfiling {
		profileTypes = []profiler.ProfileType{
			profiler.CPUProfile,
			profiler.HeapProfile,
			profiler.GoroutineProfile,
			profiler.MutexProfile,
			profiler.BlockProfile,
		}
	} else {
		profileTypes = []profiler.ProfileType{profiler.CPUProfile}
	}

	err := profiler.Start(
		profiler.WithEnv(cfg.Env),
		profiler.WithUDS(cfg.APM),
		profiler.WithService(cfg.Service),
		profiler.WithVersion(cfg.ServiceVersion),
		profiler.WithProfileTypes(profileTypes...),
	)
	if err != nil {
		logger.Errorf("Failed to start Datadog profiler: %v", err)
	}
}

// GracefulDatadogShutdown of executed parallel processes
func GracefulDatadogShutdown() {
	defer tracer.Stop()
	defer profiler.Stop()
}
