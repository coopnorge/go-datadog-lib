package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsDataDogConfigValid(t *testing.T) {
	cfg := DatadogConfig{}

	assert.False(t, IsDataDogConfigValid(cfg))

	cfg.Env = "dev"
	assert.False(t, IsDataDogConfigValid(cfg))

	cfg.Service = "Lib"
	assert.False(t, IsDataDogConfigValid(cfg))

	cfg.ServiceVersion = "v1"
	assert.False(t, IsDataDogConfigValid(cfg))

	cfg.DSD = "/tmp"
	assert.False(t, IsDataDogConfigValid(cfg))

	cfg.DSD = "unix///tmp"
	cfg.APM = "unix///tmp"
	assert.False(t, IsDataDogConfigValid(cfg))

	cfg.APM = "/tmp"
	assert.True(t, IsDataDogConfigValid(cfg))
}
