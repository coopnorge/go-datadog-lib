package internal

import (
	"os"
	"strconv"
)

const (
	// DatadogDisable is the environment variable key for whether to disable the Datadog integration.
	DatadogDisable = "DD_DISABLE"
	// DatadogEnableTracing is the environment variable key for whether to enable tracing or not.
	DatadogEnableTracing = "DD_ENABLE_TRACING"
	// DatadogEnableProfiling is the environment variable key for whether to enable profiling or not.
	DatadogEnableProfiling = "DD_ENABLE_PROFILING"
	// DatadogEnableExtraProfiling is the environment variable key for whether to enable extra profiling or not.
	DatadogEnableExtraProfiling = "DD_ENABLE_EXTRA_PROFILING"

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
)

// IsDatadogDisabled checks if the Datadog integration is disabled. The
// environment variable DD_DISABLE is checked. If the variable is missing or
// cannot be parsed to a bool the Datadog integration is assumed to be enabled.
func IsDatadogDisabled() bool {
	valStr := os.Getenv(DatadogDisable)
	val, err := strconv.ParseBool(valStr)
	if err != nil {
		return false
	}
	return val
}
