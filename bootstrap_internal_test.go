package coopdatadog

import (
	"os"
	"testing"

	"github.com/coopnorge/go-datadog-lib/v2/internal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

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
