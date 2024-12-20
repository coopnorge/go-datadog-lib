package internal

import (
	"fmt"
	"os"
	"strconv"
)

const (
	// DatadogDisable is the environment variable key for whether to disable the Datadog integration.
	DatadogDisable = "DD_DISABLE"
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
	return GetBool(DatadogDisable, false)
}

// VerifyEnvVarsSet checks if the provided environmental variables are defined
func VerifyEnvVarsSet(keys ...string) error {
	for _, key := range keys {
		val, ok := os.LookupEnv(key)
		if !ok || val == "" {
			return fmt.Errorf("required environmental variable not set: %q", key)
		}
	}
	return nil
}

// GetBool returns the boolean value of the environmental variable, if the key
// is not set or parsing fails the fallback value is returned.
func GetBool(key string, fallback bool) bool {
	valStr, ok := os.LookupEnv(key)
	if !ok {
		return fallback
	}
	val, err := strconv.ParseBool(valStr)
	if err != nil {
		return fallback
	}
	return val
}
