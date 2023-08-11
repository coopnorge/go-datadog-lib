package config

import (
	"strconv"
	"strings"
)

type (
	// DatadogParameters for connection and configuring background process to send information to Datadog Agent
	DatadogParameters interface {
		// GetEnv where application is executed, dev, production, staging etc
		GetEnv() string
		// GetService how must be service called and displayed in Datadog system
		GetService() string
		// GetServiceVersion depends on system, can be Git Tag or API version
		GetServiceVersion() string
		// GetDsdEndpoint Socket path or URL for DD StatsD
		GetDsdEndpoint() string
		// GetApmEndpoint Socket path or URL for APM and profiler
		GetApmEndpoint() string
		// IsExtraProfilingEnabled flag enables more optional profilers not recommended for production.
		IsExtraProfilingEnabled() bool
		// IsDataDogConfigValid method to verify if configuration values are correct
		IsDataDogConfigValid() bool
	}
	// DatadogConfig that required to connect to Datadog Agent
	DatadogConfig struct {
		// Env where application is executed, dev, production, staging etc
		Env string `mapstructure:"dd_env" json:"dd_env,omitempty"`
		// Service how must be service called and displayed in Datadog system
		Service string `mapstructure:"dd_service" json:"dd_service,omitempty"`
		// ServiceVersion depends on system, can be Git Tag or API version
		ServiceVersion string `mapstructure:"dd_version" json:"dd_version,omitempty"`
		// DSD Socket path for DD StatsD, important to have unix prefix for that value, example: unix:///var/run/dd/dsd.socket
		DSD string `mapstructure:"dd_dogstatsd_url" json:"dd_dogstatsd_url,omitempty"`
		// APM Socket path for APM and profiler, unix prefix not needed, example: /var/run/dd/apm.socket
		APM string `mapstructure:"dd_trace_agent_url" json:"dd_trace_agent_url,omitempty"`
		// EnableExtraProfiling bool flag enables more optional profilers not recommended for production
		EnableExtraProfiling string `mapstructure:"dd_enable_extra_profiling" json:"dd_enable_extra_profiling,omitempty"`
	}
)

// IsDataDogConfigValid method to verify if configuration values are correct
func (d DatadogConfig) IsDataDogConfigValid() bool {
	if d.Env == "" {
		return false
	}
	if d.Service == "" {
		return false
	}
	if d.ServiceVersion == "" {
		return false
	}

	// DSD or APM must be configured`
	if d.DSD == "" && d.APM == "" {
		return false
	}

	return true
}

// GetEnv where application is executed, dev, production, staging etc
func (d DatadogConfig) GetEnv() string {
	return d.Env
}

// GetService how must be service called and displayed in Datadog system
func (d DatadogConfig) GetService() string {
	return d.Service
}

// GetServiceVersion depends on system, can be Git Tag or API version
func (d DatadogConfig) GetServiceVersion() string {
	return d.ServiceVersion
}

// GetDsdEndpoint Socket path or URL for DD StatsD
// for socket important to have unix prefix for that value, example: unix:///var/run/dd/dsd.socket
func (d DatadogConfig) GetDsdEndpoint() string {
	return d.DSD
}

// GetApmEndpoint Socket path or URL for APM and profiler
// unix prefix not needed, example: /var/run/dd/apm.socket
func (d DatadogConfig) GetApmEndpoint() string {
	return d.APM
}

// IsExtraProfilingEnabled return true if profilers not recommended for production are enabled.
func (d DatadogConfig) IsExtraProfilingEnabled() bool {
	isEnabled, failedToConvert := strconv.ParseBool(strings.ToLower(strings.TrimSpace(d.EnableExtraProfiling)))
	if failedToConvert != nil {
		return false
	}

	return isEnabled
}
