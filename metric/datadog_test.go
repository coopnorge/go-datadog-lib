package metrics

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

	m := NewDatadogMetrics(&cfg)

	assert.True(t, m.defaultMetricsTags[0] == "environment:dev")
	assert.True(t, m.defaultMetricsTags[1] == "service:unit-test")
	assert.True(t, m.defaultMetricsTags[2] == "version:VUnit")

	assert.True(t, m.servicePrefix == "unit_test")

	assert.Nil(t, m.GetClient())
	assert.ElementsMatch(t, m.defaultMetricsTags, m.GetDefaultTags())
	assert.Equal(t, m.GetServiceNamePrefix(), "unit_test")
}
