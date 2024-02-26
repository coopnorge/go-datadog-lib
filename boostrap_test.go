package coopdatadog

import (
	"testing"

	"github.com/coopnorge/go-datadog-lib/v2/config"
	"github.com/stretchr/testify/assert"
)

func TestDatadog(t *testing.T) {
	ddCfg := new(config.DatadogConfig)

	err := StartDatadog(ddCfg, ConnectionTypeHTTP)
	assert.NotNil(t, err)

	ddCfg = &config.DatadogConfig{
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
