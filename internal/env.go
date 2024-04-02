package internal

import (
	"os"
	"strconv"
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
	// DatadogEnable is the environment variable key for whether to enable the Datadog integration.
	DatadogEnable = "DD_ENABLE"
)

// IsDatadogConfigured checks some common environment-variables to determine if
// the service is configured to use Datadog.
func IsDatadogConfigured() bool {
	if !IsDatadogEnabled() {
		return false
	}
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

// IsDatadogEnabled checks if the Datadog integration is enabled. The
// environment variable DD_ENABLE is checked. If the variable is missing the
// Datadog integration is assumed to be enabled.
func IsDatadogEnabled() bool {
	valStr := os.Getenv(DatadogEnable)
	val, err := strconv.ParseBool(valStr)
	if err != nil {
		return true
	}
	return val
}
