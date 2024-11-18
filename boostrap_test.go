package coopdatadog

import (
	"testing"

	"github.com/coopnorge/go-datadog-lib/v2/config"
	"github.com/coopnorge/go-datadog-lib/v2/internal"
	"github.com/stretchr/testify/assert"
)

func TestDatadog(t *testing.T) {
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
