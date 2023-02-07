package metric

import (
	"testing"

	"github.com/coopnorge/go-datadog-lib/config"

	"github.com/stretchr/testify/assert"
)

func TestNewDatadogMetrics(t *testing.T) {
	cfg := config.DatadogConfig{
		Env:            "dev",
		Service:        "unit-test",
		ServiceVersion: "VUnit",
	}

	m, err := NewDatadogMetrics(&cfg, "myUnitTest")
	assert.NotNil(t, err)
	assert.Nil(t, m)
}
