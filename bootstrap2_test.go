package coopdatadog

import (
	"context"
	"testing"

	"github.com/coopnorge/go-datadog-lib/v2/internal"
	"github.com/stretchr/testify/assert"
)

func TestBootstrapDatadogDisabled(t *testing.T) {
	t.Setenv(internal.DatadogDisable, "true")

	orgCtx := context.Background()

	ctx, cancel, err := Start(orgCtx)

	assert.NoError(t, err)
	assert.Equal(t, orgCtx, ctx)
	assert.NotNil(t, cancel)
}

func TestBootstrap(t *testing.T) {
	t.Setenv(internal.DatadogDisable, "false")
	t.Setenv(internal.DatadogAPMEndpoint, "/tmp")
	t.Setenv(internal.DatadogDSDEndpoint, "unix:///tmp/")
	t.Setenv(internal.DatadogEnvironment, "unittest")
	t.Setenv(internal.DatadogService, "go-datadog-lib-unit-test")
	t.Setenv(internal.DatadogVersion, "42345kjh435")

	orgCtx := context.Background()

	ctx, cancel, err := Start(orgCtx)
	assert.NoError(t, err)
	assert.NotEqual(t, orgCtx, ctx)
	assert.NotNil(t, cancel)

	cancel()
}

func TestBootstrapMissingEnvVar(t *testing.T) {
	t.Setenv(internal.DatadogDisable, "false")

	orgCtx := context.Background()

	ctx, cancel, err := Start(orgCtx)
	assert.ErrorContains(t, err, "required environmental variable not set: ")
	assert.Equal(t, orgCtx, ctx)
	assert.NotNil(t, cancel)
}
