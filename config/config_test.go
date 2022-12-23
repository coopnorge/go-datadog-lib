package config

import (
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
    assert. True(t, cfg.IsDataDogConfigValid())

    cfg.APM = "unix///tmp"
    cfg.DSD = ""
    assert. True(t, cfg.IsDataDogConfigValid())

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
}
