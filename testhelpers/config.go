package testhelpers

import (
	"strconv"
	"testing"

	"github.com/coopnorge/go-datadog-lib/v2/config"
	"github.com/coopnorge/go-datadog-lib/v2/internal"
	"github.com/stretchr/testify/assert"
)

// ConfigureDatadog prepares the environment for running tests
func ConfigureDatadog(t *testing.T) {
	t.Setenv(internal.DatadogEnable, strconv.FormatBool(true))
	t.Setenv(internal.DatadogEnvironment, "unittest")
	t.Setenv(internal.DatadogService, "unittest-service")
	t.Setenv(internal.DatadogVersion, "v0.0.0")
	t.Setenv(internal.DatadogAPMEndpoint, "/dev/null")
	t.Setenv(internal.DatadogDSDEndpoint, "unix:///dev/null")
	cfg := config.LoadDatadogConfigFromEnvVars()
	err := cfg.Validate()
	assert.NoError(t, err)
}
