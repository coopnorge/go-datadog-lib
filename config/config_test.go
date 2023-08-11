package config

import (
	"strconv"
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

	testParsingStringToBool := []struct {
		input    string
		expected bool
		valid    bool
	}{
		{"true", true, true},
		{"t", true, true},
		{"1", true, true},
		{"y", false, false},
		{"yes", false, false},
		{"false", false, true},
		{"f", false, true},
		{"0", false, true},
		{"n", false, false},
		{"no", false, false},
		{"invalid", false, false},
	}

	for _, tt := range testParsingStringToBool {
		t.Run(tt.input, func(t *testing.T) {
			expectedCfg.EnableExtraProfiling = tt.input

			isBool, parseFailed := strconv.ParseBool(expectedCfg.EnableExtraProfiling)
			if parseFailed != nil && tt.valid {
				t.Fatalf("Expected no error but got '%v'", parseFailed)
			}

			if parseFailed == nil && !tt.valid {
				t.Fatalf("Expected an error for input '%v' but got none", tt.input)
			}

			if tt.valid && isBool != tt.expected {
				t.Fatalf("Expected '%v' but got '%v' for input '%v'", tt.expected, isBool, tt.input)
			}
		})
	}
}
