package internal

import (
	"os"
	"strings"
)

const ExperimentalTracingEnabled = "DD_EXPERIMENTAL_TRACING_ENABLED"

func IsExperimentalTracingEnabled() bool {
	if strings.ToLower(os.Getenv(ExperimentalTracingEnabled)) == "true" {
		return true
	}
	return false
}
