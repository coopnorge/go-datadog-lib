package config

import "strings"

type (
	// DatadogConfig that required to connect to Datadog Agent
	DatadogConfig struct {
		// Env where application is executed, dev, production, staging etc
		Env string `mapstructure:"dd_env"`
		// Service how must be service called and displayed in Datadog system
		Service string `mapstructure:"dd_service"`
		// ServiceVersion depends on system, can be Git Tag or API version
		ServiceVersion string `mapstructure:"dd_version"`
		// DSD Socket path for DD StatsD, important to have unix prefix for that value, example: unix:///var/run/dd/dsd.socket
		DSD string `mapstructure:"dd_dogstatsd_url"`
		// APM Socket path for apm and profiler, unix prefix not needed, example: /var/run/dd/apm.socket
		APM string `mapstructure:"dd_trace_agent_url"`
	}
)

// IsDataDogConfigValid with given values
func IsDataDogConfigValid(cfg DatadogConfig) bool {
	if cfg.Env == "" {
		return false
	}
	if cfg.Service == "" {
		return false
	}
	if cfg.ServiceVersion == "" {
		return false
	}

	// Check socket paths
	if cfg.DSD != "" && strings.Contains(cfg.DSD, "unix") {
		return true
	}
	if cfg.APM != "" && !strings.Contains(cfg.DSD, "unix") {
		return true
	}

	// DSD or APM must be configured
	if cfg.DSD == "" && cfg.APM == "" {
		return false
	}

	return true
}
