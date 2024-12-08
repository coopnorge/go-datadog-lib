package coopdatadog

import (
	"context"
	"testing"

	"github.com/coopnorge/go-datadog-lib/v2/internal"
	"github.com/stretchr/testify/assert"
)

func TestBootstrapDatadogDisabled(t *testing.T) {
	t.Setenv(internal.DatadogDisable, "true")

	stop, err := Start(context.Background())
	defer func() {
		err := stop()
		assert.NoError(t, err)
	}()

	assert.NoError(t, err)
	assert.NotNil(t, stop)
}

func TestBootstrap(t *testing.T) {
	t.Setenv(internal.DatadogDisable, "false")
	t.Setenv(internal.DatadogAPMEndpoint, "/tmp")
	t.Setenv(internal.DatadogDSDEndpoint, "unix:///tmp/")
	t.Setenv(internal.DatadogEnvironment, "unittest")
	t.Setenv(internal.DatadogService, "go-datadog-lib-unit-test")
	t.Setenv(internal.DatadogVersion, "42345kjh435")

	stop, err := Start(context.Background())
	defer func() {
		err := stop()
		assert.NoError(t, err)
	}()

	assert.NoError(t, err)
	assert.NotNil(t, stop)
}

func TestBootstrapMissingEnvVar(t *testing.T) {
	t.Setenv(internal.DatadogDisable, "false")

	stop, err := Start(context.Background())
	defer func() {
		err := stop()
		assert.NoError(t, err)
	}()
	assert.ErrorContains(t, err, "required environmental variable not set: ")
	assert.NotNil(t, stop)
}
