package go_datadog_lib

import (
	"testing"

	"github.com/coopnorge/go-datadog-lib/config"
)

func TestDatadog(t *testing.T) {
	ddCfg := config.DatadogConfig{}

	StartDatadog(ddCfg, false)

	ddCfg = config.DatadogConfig{
		Env:            "local",
		Service:        "Test-Go-Datadog-lib",
		ServiceVersion: "na",
		DSD:            "unix:///tmp/",
		APM:            "/tmp",
	}

	StartDatadog(ddCfg, true)
	GracefulDatadogShutdown()
}
