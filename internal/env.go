package internal

import (
	"os"
)

const (
	// DatadogEnvironment is the environment variable key determining the Datadog Environment to use.
	DatadogEnvironment = "DD_ENV"
	// DatadogService is the environment variable key for the name of the current service.
	DatadogService = "DD_SERVICE"
	// DatadogVersion is the environment variable key for the version of the current service.
	DatadogVersion = "DD_VERSION"
	// DatadogDSDEndpoint is the environment variable key for the URL to StatsD.
	DatadogDSDEndpoint = "DD_DOGSTATSD_URL"
	// DatadogAPMEndpoint is the environment variable key for the URL to APM.
	DatadogAPMEndpoint = "DD_TRACE_AGENT_URL"
	// DatadogEnableExtraProfiling is the environment variable key for whether to enable extra profiling or not.
	DatadogEnableExtraProfiling = "DD_ENABLE_EXTRA_PROFILING"
)

// IsDatadogConfigured checks some common environment-variables to determine if the service is configured to use Datadog.
func IsDatadogConfigured() bool {
	if val := os.Getenv(DatadogEnvironment); val != "" {
		return true
	}
	if val := os.Getenv(DatadogService); val != "" {
		return true
	}
	if val := os.Getenv(DatadogVersion); val != "" {
		return true
	}
	return false
}
