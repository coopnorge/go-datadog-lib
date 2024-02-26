package config

import (
	"encoding/json"
	"github.com/spf13/viper"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsDataDogConfigValid(t *testing.T) {
	cfg := DatadogConfig{}

	assert.False(t, cfg.IsDataDogConfigValid())

	cfg.Env = "dev"
	assert.False(t, cfg.IsDataDogConfigValid())

	cfg.Service = "Lib"
	assert.False(t, cfg.IsDataDogConfigValid())

	cfg.ServiceVersion = "v1"
	assert.False(t, cfg.IsDataDogConfigValid())

	cfg.APM = ""
	cfg.DSD = ""
	assert.False(t, cfg.IsDataDogConfigValid())

	cfg.APM = ""
	cfg.DSD = "unix///tmp"
	assert.True(t, cfg.IsDataDogConfigValid())

	cfg.APM = "unix///tmp"
	cfg.DSD = ""
	assert.True(t, cfg.IsDataDogConfigValid())

	cfg.DSD = "unix///tmp"
	cfg.APM = "unix///tmp"
	assert.True(t, cfg.IsDataDogConfigValid())
}

func TestValidate(t *testing.T) {
	cfg := DatadogConfig{}

	assert.Error(t, cfg.Validate())

	cfg.Env = "dev"
	assert.Error(t, cfg.Validate())

	cfg.Service = "Lib"
	assert.Error(t, cfg.Validate())

	cfg.ServiceVersion = "v1"
	assert.Error(t, cfg.Validate())

	cfg.APM = ""
	cfg.DSD = ""
	assert.Error(t, cfg.Validate())

	cfg.APM = ""
	cfg.DSD = "unix///tmp"
	assert.Nil(t, cfg.Validate())

	cfg.APM = "unix///tmp"
	cfg.DSD = ""
	assert.Nil(t, cfg.Validate())

	cfg.DSD = "unix///tmp"
	cfg.APM = "unix///tmp"
	assert.Nil(t, cfg.Validate())
}

func TestConfigGetters(t *testing.T) {
	expectedCfg := DatadogConfig{
		Env:            "unit",
		Service:        "Service",
		ServiceVersion: "ServiceVersion",
		DSD:            "DSD:",
		APM:            "APM",
	}

	assert.Equal(t, expectedCfg.Env, expectedCfg.GetEnv())
	assert.Equal(t, expectedCfg.Service, expectedCfg.GetService())
	assert.Equal(t, expectedCfg.ServiceVersion, expectedCfg.GetServiceVersion())
	assert.Equal(t, expectedCfg.DSD, expectedCfg.GetDsdEndpoint())
	assert.Equal(t, expectedCfg.APM, expectedCfg.GetApmEndpoint())
}

func TestUnmarshalDatadogConfigWithVanillaJsonTag(t *testing.T) {
	const fixtureConfig = `{"dd_env": "unit_test", "dd_service": "go_test", "dd_service_version": "unit", "dd_dsd": "gl", "dd_apm": "hf", "dd_enable_extra_profiling": "true"}`
	var dogCfg DatadogConfig
	err := json.Unmarshal([]byte(fixtureConfig), &dogCfg)

	assert.NoError(t, err)
	assert.NotEmpty(t, dogCfg)
	assert.True(t, dogCfg.Env == "unit_test", "Env expected to be => unit_test")
	assert.True(t, dogCfg.Service == "go_test", "Env expected to be => go_test")
	assert.True(t, dogCfg.ServiceVersion == "unit", "Env expected to be => unit")
	assert.True(t, dogCfg.DSD == "gl", "Env expected to be => gl")
	assert.True(t, dogCfg.APM == "hf", "Env expected to be => hf")
	assert.True(t, dogCfg.EnableExtraProfiling, "EnableExtraProfiling expected to be => true")
}

func TestUnmarshalDatadogConfigWithMapstructureAndViper(t *testing.T) {
	testCases := []struct {
		name     string
		envField string
		envValue string
	}{
		// mapstructure json tag
		{
			name:     "mapstructure - dd_env",
			envField: "dd_env",
			envValue: "unit_test",
		},
		{
			name:     "mapstructure - dd_service",
			envField: "dd_service",
			envValue: "go_test",
		},
		{
			name:     "mapstructure - dd_version",
			envField: "dd_version",
			envValue: "unit",
		},
		{
			name:     "mapstructure - dd_dogstatsd_url",
			envField: "dd_dogstatsd_url",
			envValue: "gl",
		},
		{
			name:     "mapstructure - dd_trace_agent_url",
			envField: "dd_trace_agent_url",
			envValue: "hf",
		},
		{
			name:     "mapstructure - dd_enable_extra_profiling",
			envField: "dd_enable_extra_profiling",
			envValue: "true",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			envName := strings.ToUpper(tc.envField)
			_ = os.Setenv(envName, tc.envValue)

			viper.AutomaticEnv()

			assert.Equal(t, tc.envValue, viper.GetString(tc.envField))
		})
	}

}
