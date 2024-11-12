package coopdatadog

import (
	"testing"

	"github.com/coopnorge/go-datadog-lib/v2/config"
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

func TestSetConnectionType(t *testing.T) {
	ddCfg := config.DatadogConfig{
		Env:                  "local",
		Service:              "Test-Go-Datadog-lib",
		ServiceVersion:       "na",
		DSD:                  "unix:///tmp/",
		APM:                  "http://localhost:3899",
		EnableExtraProfiling: true,
	}

	// If not auto it should just pass through
	connectionType, err := setConnectionType(ddCfg, ConnectionTypeSocket)
	assert.NoError(t, err)
	assert.Equal(t, ConnectionTypeSocket, connectionType)

	// Auto should detect
	connectionType, err = setConnectionType(ddCfg, ConnectionTypeAuto)
	assert.NoError(t, err)
	assert.Equal(t, ConnectionTypeHTTP, connectionType)

	// Fail on unable to detect when auto
	ddCfg.APM = "tmp"
	connectionType, err = setConnectionType(ddCfg, ConnectionTypeAuto)
	assert.Error(t, err)
	assert.Equal(t, ConnectionTypeAuto, connectionType)
}
