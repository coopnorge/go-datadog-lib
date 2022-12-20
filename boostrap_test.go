package coopdatadog

import (
	"testing"

	"github.com/coopnorge/go-datadog-lib/config"
	"github.com/stretchr/testify/assert"
)

func TestDatadog(t *testing.T) {
	ddCfg := config.DatadogConfig{}

	err := StartDatadog(ddCfg, false, false)
	assert.NotNil(t, err)

	ddCfg = config.DatadogConfig{
		Env:            "local",
		Service:        "Test-Go-Datadog-lib",
		ServiceVersion: "na",
		DSD:            "unix:///tmp/",
		APM:            "/tmp",
	}

	err = StartDatadog(ddCfg, true, true)
    assert.Nil(t, err)

	GracefulDatadogShutdown()
}
