package coopdatadog

import (
	"context"
	"testing"

	"github.com/coopnorge/go-datadog-lib/v2/internal"
	"github.com/coopnorge/go-datadog-lib/v2/internal/testhelpers"
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
	testhelpers.ConfigureDatadog(t)

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
