package metric_test

import (
	"testing"

	"github.com/coopnorge/go-datadog-lib/v2/config"
	"github.com/coopnorge/go-datadog-lib/v2/metric"

	"github.com/stretchr/testify/assert"
)

func TestNewDatadogMetrics(t *testing.T) {
	cfg := config.DatadogConfig{
		Env:            "dev",
		Service:        "unit-test",
		ServiceVersion: "VUnit",
	}

	m, err := metric.NewDatadogMetrics(&cfg, "myUnitTest")
	assert.NotNil(t, err)
	assert.Nil(t, m)
}
