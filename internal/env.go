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
