package coopdatadog

import (
	"os"
	"testing"

	"github.com/coopnorge/go-datadog-lib/v2/config"
	"github.com/coopnorge/go-datadog-lib/v2/internal"
	"github.com/stretchr/testify/assert"
)

func TestDatadog(t *testing.T) {
	t.Cleanup(func() {
		os.Unsetenv(internal.DatadogEnvironment)
		os.Unsetenv(internal.DatadogService)
		os.Unsetenv(internal.DatadogVersion)
		os.Unsetenv(internal.DatadogDSDEndpoint)
		os.Unsetenv(internal.DatadogAPMEndpoint)
	})

	ddCfg := config.DatadogConfig{}

	err := StartDatadog(ddCfg, ConnectionTypeHTTP)
	assert.NotNil(t, err)

	ddCfg = config.DatadogConfig{
		Env:                  "local",
		Service:              "Test-Go-Datadog-lib",
		ServiceVersion:       "na",
		DSD:                  "unix:///tmp/",
		APM:                  "/tmp",
		EnableExtraProfiling: true,
	}

	err = StartDatadog(ddCfg, ConnectionTypeSocket)
	assert.Nil(t, err)

	GracefulDatadogShutdown()
}

func TestValidateConnectionType(t *testing.T) {
	testCases := map[string]struct {
		envVal    string
		connType  ConnectionType
		expectErr bool
	}{
		"Socket no env":   {envVal: "", connType: ConnectionTypeSocket, expectErr: false},
		"Socket with env": {envVal: "foobar", connType: ConnectionTypeSocket, expectErr: false},
		"HTTP no env":     {envVal: "", connType: ConnectionTypeHTTP, expectErr: false},
		"HTTP with env":   {envVal: "foobar", connType: ConnectionTypeHTTP, expectErr: false},
		"Auto no env":     {envVal: "", connType: ConnectionTypeAuto, expectErr: true},
		"Auto with env":   {envVal: "foobar", connType: ConnectionTypeAuto, expectErr: false},
	}
	for k, tt := range testCases {
		t.Run(k, func(t *testing.T) {
			t.Setenv(internal.DatadogAPMEndpoint, tt.envVal)
			err := validateConnectionType(tt.connType)
			if tt.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestNormalizeLegacySocketPath(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"unix:///var/run/datadog/apm.socket", "/var/run/datadog/apm.socket"},
		{"/var/run/datadog/apm.socket", "/var/run/datadog/apm.socket"},
	}
	for _, tc := range tests {
		t.Run(tc.input, func(t *testing.T) {
			got := normalizeLegacySocketPath(tc.input)
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestNormalizeLegacyHTTPAddr(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"http://my-dd-agent:3678", "my-dd-agent:3678"},
		{"http://my-dd-agent", "my-dd-agent"},
		{"http://10.0.0.6", "10.0.0.6"},
	}
	for _, tc := range tests {
		t.Run(tc.input, func(t *testing.T) {
			got := normalizeLegacyHTTPAddr(tc.input)
			assert.Equal(t, tc.want, got)
		})
	}
}
