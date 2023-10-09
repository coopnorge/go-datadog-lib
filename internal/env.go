package internal

import (
	"os"
	"strings"
)

// ExperimentalTracingEnabled is the environment variable key determining if experimental tracing should be enabled
const ExperimentalTracingEnabled = "DD_EXPERIMENTAL_TRACING_ENABLED"

// IsExperimentalTracingEnabled checks if experimental tracing is enabled
func IsExperimentalTracingEnabled() bool {
	return strings.ToLower(os.Getenv(ExperimentalTracingEnabled)) == "true"
}
