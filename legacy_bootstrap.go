package coopdatadog

import (
	"fmt"
	"os"
	"strings"

	"github.com/coopnorge/go-datadog-lib/v2/config"
	"github.com/coopnorge/go-datadog-lib/v2/internal"
	"github.com/coopnorge/go-datadog-lib/v2/internal/log"
	"github.com/coopnorge/go-logger"

	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
	"gopkg.in/DataDog/dd-trace-go.v1/profiler"
)

// ConnectionType enum type
//
// Deprecated: Use coopdatadog.Start() instead.
type ConnectionType byte

const (
	// ConnectionTypeSocket sets the connection to Datadog to go through a UNIX socket
	//
	// Deprecated: ConnectionTypeSocket. ConnectionTypeAuto should be used.
	ConnectionTypeSocket ConnectionType = iota
	// ConnectionTypeHTTP sets the connection to Datadog to go over HTTP
	//
	// Deprecated: ConnectionTypeHTTP. ConnectionTypeAuto should be used.
	ConnectionTypeHTTP
	// ConnectionTypeAuto sets connection to HTTP or UNIX depending on supplied configuration of DD_TRACE_AGENT_URL
	//
	// Deprecated: Use coopdatadog.Start() instead.
	ConnectionTypeAuto
)

// StartDatadog parallel process to collect data for Datadog.
// connectionType flag related to Datadog connection type, it supports HTTP or socket - values will be used from config.DatadogParameters
//
// Deprecated: Use coopdatadog.Start() instead.
func StartDatadog(cfg config.DatadogParameters, connectionType ConnectionType) error {
	if internal.IsDatadogDisabled() {
		return nil
	}

	if err := cfg.Validate(); err != nil {
		return fmt.Errorf("the Datadog configuration not valid, cannot initialize Datadog services: %w", err)
	}

	l, err := log.NewLogger(log.WithGlobalLogger())
	if err != nil {
		return fmt.Errorf("failed to initialize the Datadog logger: %w", err)
	}
	ddtrace.UseLogger(l)

	compareConfigWithEnv(cfg)

	if err := validateConnectionType(connectionType); err != nil {
		return err
	}

	initTracer(cfg, connectionType)
	if initProfilerErr := initProfiler(cfg, connectionType); initProfilerErr != nil {
		return fmt.Errorf("failed to start Datadog profiler: %w", initProfilerErr)
	}

	return nil
}

func compareConfigWithEnv(cfg config.DatadogParameters) {
	envCfg := config.LoadDatadogConfigFromEnvVars()

	fields := map[string]any{}
	if cfg.GetEnv() != envCfg.GetEnv() {
		fields["input-env"] = cfg.GetEnv()
		fields["envvar-env"] = envCfg.GetEnv()
		_ = os.Setenv(internal.DatadogEnvironment, cfg.GetEnv()) //nolint:errcheck
	}
	if cfg.GetService() != envCfg.GetService() {
		fields["input-service"] = cfg.GetService()
		fields["envvar-service"] = envCfg.GetService()
		_ = os.Setenv(internal.DatadogService, cfg.GetService()) //nolint:errcheck
	}
	if cfg.GetServiceVersion() != envCfg.GetServiceVersion() {
		fields["input-service-version"] = cfg.GetServiceVersion()
		fields["envvar-service-version"] = envCfg.GetServiceVersion()
		_ = os.Setenv(internal.DatadogVersion, cfg.GetServiceVersion()) //nolint:errcheck
	}
	if cfg.GetDsdEndpoint() != envCfg.GetDsdEndpoint() {
		fields["input-dsd-url"] = cfg.GetDsdEndpoint()
		fields["envvar-dsd-url"] = envCfg.GetDsdEndpoint()
		_ = os.Setenv(internal.DatadogDSDEndpoint, cfg.GetDsdEndpoint()) //nolint:errcheck
	}
	if cfg.GetApmEndpoint() != envCfg.GetApmEndpoint() {
		fields["input-apm-url"] = cfg.GetApmEndpoint()
		fields["envvar-apm-url"] = envCfg.GetApmEndpoint()
		_ = os.Setenv(internal.DatadogAPMEndpoint, cfg.GetApmEndpoint()) //nolint:errcheck
	}

	// Note: IsExtraProfilingEnabled is internal to this library, so we won't warn or set env-var if it differs

	if len(fields) > 0 {
		logger.WithFields(fields).Warn("Supplied config does not match config from env-vars. See https://github.com/coopnorge/go-datadog-lib/issues/310")
	}
}

// GracefulDatadogShutdown of executed parallel processes
//
// Deprecated: Use the StopFunc returned coopdatadog.Start() instead.
func GracefulDatadogShutdown() {
	defer tracer.Stop()
	defer profiler.Stop()
}

func validateConnectionType(connectionType ConnectionType) error {
	if connectionType == ConnectionTypeAuto {
		// When using ConnectionTypeAuto, we offload the determining of the connection-type to the underlying library, which only reads known environment-variables.
		envVal := os.Getenv(internal.DatadogAPMEndpoint)
		if envVal == "" {
			return fmt.Errorf("to use ConnectionTypeAuto, the environment-variable %q MUST be set", internal.DatadogAPMEndpoint)
		}
	}
	return nil
}

func initTracer(cfg config.DatadogParameters, connectionType ConnectionType) {
	tracerOptions := make([]tracer.StartOption, 0, 5)
	switch connectionType {
	case ConnectionTypeSocket:
		socketPath := normalizeLegacySocketPath(cfg.GetApmEndpoint())
		tracerOptions = append(tracerOptions, tracer.WithUDS(socketPath))
	case ConnectionTypeHTTP:
		httpAddr := normalizeLegacyHTTPAddr(cfg.GetApmEndpoint())
		tracerOptions = append(tracerOptions, tracer.WithAgentAddr(httpAddr))
	case ConnectionTypeAuto:
		// Let the underlying library determine the URL from environment-variables
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

	profilerOptions := make([]profiler.Option, 0, 5)
	switch connectionType {
	case ConnectionTypeSocket:
		socketPath := normalizeLegacySocketPath(cfg.GetApmEndpoint())
		profilerOptions = append(profilerOptions, profiler.WithUDS(socketPath))
	case ConnectionTypeHTTP:
		httpAddr := normalizeLegacyHTTPAddr(cfg.GetApmEndpoint())
		profilerOptions = append(profilerOptions, profiler.WithAgentAddr(httpAddr))
	case ConnectionTypeAuto:
		// Let the underlying library determine the URL from environment-variables
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

// normalizeLegacySocketPath ensures that the socketpath is in the format the tracer.WithUDS and profiler.WithUDS expects.
func normalizeLegacySocketPath(socketPath string) string {
	// profiler.WithUDS and tracer.WithUDS expects a path without the scheme
	return strings.TrimPrefix(socketPath, "unix://")
}

// normalizeLegacySocketPath ensures that the HTTP address is in the format the tracer.WithAgentAddr and profiler.WithAgentAddr expects.
func normalizeLegacyHTTPAddr(addr string) string {
	// profiler.WithAgentAddr and tracer.WithAgentAddr expects a path without the scheme
	return strings.TrimPrefix(addr, "http://")
}
