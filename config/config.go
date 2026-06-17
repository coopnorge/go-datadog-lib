// Deprecated: Use coopdatadog.Start() instead.
package config

import (
	"errors"
	"os"
	"strconv"

	"github.com/coopnorge/go-datadog-lib/v2/internal"
)

type (
	// DatadogParameters for connection and configuring background process to send information to Datadog Agent
	//
	// Deprecated: Use coopdatadog.Start() instead.
	DatadogParameters interface {
		// GetEnv where application is executed, dev, production, staging etc
		//
		// Deprecated: Use coopdatadog.Start() instead.
		GetEnv() string
		// GetService how must be service called and displayed in Datadog system
		//
		// Deprecated: Use coopdatadog.Start() instead.
		GetService() string
		// GetServiceVersion depends on system, can be Git Tag or API version
		//
		// Deprecated: Use coopdatadog.Start() instead.
		GetServiceVersion() string
		// GetDsdEndpoint Socket path or URL for DD StatsD
		//
		// Deprecated: Use coopdatadog.Start() instead.
		GetDsdEndpoint() string
		// GetApmEndpoint Socket path or URL for APM and profiler
		//
		// Deprecated: Use coopdatadog.Start() instead.
		GetApmEndpoint() string
		// IsExtraProfilingEnabled flag enables more optional profilers not recommended for production.
		//
		// Deprecated: Use coopdatadog.Start() instead.
		IsExtraProfilingEnabled() bool
		// IsDataDogConfigValid method to verify if configuration values are correct
		//
		// Deprecated: Use Validate()
		IsDataDogConfigValid() bool
		// Validate the DatadogConfig. Returns the first error found, returns nil if
		// the configuration is good.
		//
		// Deprecated: Use coopdatadog.Start() instead.
		Validate() error
	}

	// DatadogConfig that required to connect to Datadog Agent
	//
	// Deprecated: Use coopdatadog.Start() instead.
	DatadogConfig struct {
		// Env where application is executed, dev, production, staging etc
		//
		// Deprecated: Use coopdatadog.Start() instead.
		Env string `mapstructure:"dd_env" json:"dd_env,omitempty"`
		// Service how must be service called and displayed in Datadog system
		//
		// Deprecated: Use coopdatadog.Start() instead.
		Service string `mapstructure:"dd_service" json:"dd_service,omitempty"`
		// ServiceVersion depends on system, can be Git Tag or API version
		//
		// Deprecated: Use coopdatadog.Start() instead.
		ServiceVersion string `mapstructure:"dd_version" json:"dd_service_version,omitempty"`
		// DSD Socket path for DD StatsD, important to have unix prefix for that value, example: unix:///var/run/dd/dsd.socket
		//
		// Deprecated: Use coopdatadog.Start() instead.
		DSD string `mapstructure:"dd_dogstatsd_url" json:"dd_dsd,omitempty"`
		// APM Socket path for apm and profiler, unix prefix recommended, but not required, example: unix:///var/run/dd/apm.socket
		//
		// Deprecated: Use coopdatadog.Start() instead.
		APM string `mapstructure:"dd_trace_agent_url" json:"dd_apm,omitempty"`
		// EnableExtraProfiling flag enables more optional profilers not recommended for production.
		//
		// Deprecated: Use coopdatadog.Start() instead.
		EnableExtraProfiling bool `mapstructure:"dd_enable_extra_profiling" json:"dd_enable_extra_profiling,omitempty"`
	}
)

// IsDataDogConfigValid method to verify if configuration values are correct
//
// Deprecated: Use Validate()
func (d DatadogConfig) IsDataDogConfigValid() bool {
	if err := d.Validate(); err != nil {
		return false
	}

	return true
}

// Validate the DatadogConfig. Returns the first error found, returns nil if
// the configuration is good.
//
// Deprecated: Use coopdatadog.Start() instead.
func (d DatadogConfig) Validate() error {
	if d.Env == "" {
		return errors.New("DD_ENV must be defined")
	}
	if d.Service == "" {
		return errors.New("DD_SERVICE must be defined")
	}
	if d.ServiceVersion == "" {
		return errors.New("DD_VERSION must be defined")
	}

	if d.DSD == "" && d.APM == "" {
		return errors.New("DD_DOGSTATSD_URL and/or DD_TRACE_AGENT_URL must be defined")
	}

	return nil
}

// GetEnv where application is executed, dev, production, staging etc
//
// Deprecated: Use coopdatadog.Start() instead.
func (d DatadogConfig) GetEnv() string {
	return d.Env
}

// GetService how must be service called and displayed in Datadog system
//
// Deprecated: Use coopdatadog.Start() instead.
func (d DatadogConfig) GetService() string {
	return d.Service
}

// GetServiceVersion depends on system, can be Git Tag or API version
//
// Deprecated: Use coopdatadog.Start() instead.
func (d DatadogConfig) GetServiceVersion() string {
	return d.ServiceVersion
}

// GetDsdEndpoint Socket path or URL for DD StatsD
// For unix sockets, the unix-scheme prefix is required.
// Example: unix:///var/run/dd/dsd.socket
//
// Deprecated: Use coopdatadog.Start() instead.
func (d DatadogConfig) GetDsdEndpoint() string {
	return d.DSD
}

// GetApmEndpoint Socket path or URL for APM and profiler
// For unix sockets, the unix-scheme prefix is not needed, but it is recommended to include it.
// Example: unix:///var/run/dd/apm.socket
// Example: http://my-agent:1234
//
// Deprecated: Use coopdatadog.Start() instead.
func (d DatadogConfig) GetApmEndpoint() string {
	return d.APM
}

// IsExtraProfilingEnabled return true if profilers not recommended for production are enabled.
//
// Deprecated: Use coopdatadog.Start() instead.
func (d DatadogConfig) IsExtraProfilingEnabled() bool {
	return d.EnableExtraProfiling
}

// LoadDatadogConfigFromEnvVars loads a new DatadogConfig from known environment-variables.
//
// Deprecated: Use coopdatadog.Start() instead.
func LoadDatadogConfigFromEnvVars() DatadogConfig {
	return DatadogConfig{
		Env:                  os.Getenv(internal.DatadogEnvironment),
		Service:              os.Getenv(internal.DatadogService),
		ServiceVersion:       os.Getenv(internal.DatadogVersion),
		DSD:                  os.Getenv(internal.DatadogDSDEndpoint),
		APM:                  os.Getenv(internal.DatadogAPMEndpoint),
		EnableExtraProfiling: getBoolEnv(internal.DatadogEnableExtraProfiling),
	}
}

func getBoolEnv(key string) bool {
	valStr := os.Getenv(key)
	val, err := strconv.ParseBool(valStr)
	if err != nil {
		return false
	}
	return val
}
