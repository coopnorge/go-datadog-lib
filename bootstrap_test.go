package coopdatadog

import (
	"context"
	"os"
	"testing"

	"github.com/coopnorge/go-datadog-lib/v2/internal"
	"github.com/coopnorge/go-datadog-lib/v2/internal/testhelpers"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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

func TestNormalizeAPMEndpointEnv(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"unix:///var/run/datadog/apm.socket", "unix:///var/run/datadog/apm.socket"}, // Do not change a socket path that is already prefixed
		{"/var/run/datadog/apm.socket", "unix:///var/run/datadog/apm.socket"},        // Prefix an assumed unix socket path
		{"http://my-dd-agent:3678", "http://my-dd-agent:3678"},                       // Do not change HTTP-addresses
		{"http://my-dd-agent", "http://my-dd-agent"},                                 // Do not change HTTP-addresses
		{"http://10.0.0.6", "http://10.0.0.6"},                                       // Do not change HTTP-addresses
	}
	for _, tc := range tests {
		t.Run(tc.input, func(t *testing.T) {
			t.Setenv(internal.DatadogAPMEndpoint, tc.input)
			err := normalizeDatadogEnvVars()
			require.NoError(t, err)
			got := os.Getenv(internal.DatadogAPMEndpoint)
			assert.Equal(t, tc.want, got)
		})
	}
}
