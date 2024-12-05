package coopdatadog

import (
	"testing"

	"github.com/coopnorge/go-datadog-lib/v2/internal"
	"github.com/stretchr/testify/assert"
)

func TestBootstrapDatadogDisabled(t *testing.T) {
	t.Setenv(internal.DatadogDisable, "true")

	cancel, err := Start()

	assert.NoError(t, err)
	assert.NotNil(t, cancel)
}

func TestBootstrap(t *testing.T) {
	t.Setenv(internal.DatadogDisable, "false")
	t.Setenv(internal.DatadogAPMEndpoint, "/tmp")
	t.Setenv(internal.DatadogDSDEndpoint, "unix:///tmp/")
	t.Setenv(internal.DatadogEnvironment, "unittest")
	t.Setenv(internal.DatadogService, "go-datadog-lib-unit-test")
	t.Setenv(internal.DatadogVersion, "42345kjh435")

	cancel, err := Start()
	assert.NoError(t, err)
	assert.NotNil(t, cancel)

	err = cancel()
	assert.NoError(t, err)
}

func TestBootstrapMissingEnvVar(t *testing.T) {
	t.Setenv(internal.DatadogDisable, "false")

	cancel, err := Start()
	assert.ErrorContains(t, err, "required environmental variable not set: ")
	assert.NotNil(t, cancel)
}
