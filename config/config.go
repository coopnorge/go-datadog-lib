// Package config holds the Datadog configuration and reads it from the
// environment.
//
// Deprecated: The whole module is deprecated. There is no config struct. Start
// in github.com/coopnorge/go-telemetry-lib reads the standard OTEL_*
// environment variables. See the root package documentation for the migration
// guide and ADR.
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
	// Deprecated: The whole module is deprecated. There is no config struct. Start in
	// github.com/coopnorge/go-telemetry-lib reads the standard OTEL_* environment variables.
	DatadogParameters interface {
		// GetEnv where application is executed, dev, production, staging etc
		//
		// Deprecated: The whole module is deprecated. There is no config struct. Start in
		// github.com/coopnorge/go-telemetry-lib reads the standard OTEL_* environment variables.
		GetEnv() string
		// GetService how must be service called and displayed in Datadog system
		//
		// Deprecated: The whole module is deprecated. There is no config struct. Start in
		// github.com/coopnorge/go-telemetry-lib reads the standard OTEL_* environment variables.
		GetService() string
		// GetServiceVersion depends on system, can be Git Tag or API version
		//
		// Deprecated: The whole module is deprecated. There is no config struct. Start in
		// github.com/coopnorge/go-telemetry-lib reads the standard OTEL_* environment variables.
		GetServiceVersion() string
		// GetDsdEndpoint Socket path or URL for DD StatsD
		//
		// Deprecated: The whole module is deprecated. There is no config struct. Start in
		// github.com/coopnorge/go-telemetry-lib reads the standard OTEL_* environment variables.
		GetDsdEndpoint() string
		// GetApmEndpoint Socket path or URL for APM and profiler
		//
		// Deprecated: The whole module is deprecated. There is no config struct. Start in
		// github.com/coopnorge/go-telemetry-lib reads the standard OTEL_* environment variables.
		GetApmEndpoint() string
		// IsExtraProfilingEnabled flag enables more optional profilers not recommended for production.
		//
		// Deprecated: The whole module is deprecated. There is no config struct. Start in
		// github.com/coopnorge/go-telemetry-lib reads the standard OTEL_* environment variables.
		IsExtraProfilingEnabled() bool
		// IsDataDogConfigValid method to verify if configuration values are correct
		//
		// Deprecated: Use Validate().
		// The whole module is deprecated. There is no config struct. Start in
		// github.com/coopnorge/go-telemetry-lib reads the standard OTEL_* environment variables.
		IsDataDogConfigValid() bool
		// Validate the DatadogConfig. Returns the first error found, returns nil if
		// the configuration is good.
		//
		// Deprecated: The whole module is deprecated. There is no config struct. Start in
		// github.com/coopnorge/go-telemetry-lib reads the standard OTEL_* environment variables.
		Validate() error
	}

	// DatadogConfig that required to connect to Datadog Agent
	//
	// Deprecated: The whole module is deprecated. There is no config struct. Start in
	// github.com/coopnorge/go-telemetry-lib reads the standard OTEL_* environment variables.
	DatadogConfig struct {
		// Env where application is executed, dev, production, staging etc
		//
		// Deprecated: The whole module is deprecated. There is no config struct. Start in
		// github.com/coopnorge/go-telemetry-lib reads the standard OTEL_* environment variables.
		Env string `mapstructure:"dd_env" json:"dd_env,omitempty"`
		// Service how must be service called and displayed in Datadog system
		//
		// Deprecated: The whole module is deprecated. There is no config struct. Start in
		// github.com/coopnorge/go-telemetry-lib reads the standard OTEL_* environment variables.
		Service string `mapstructure:"dd_service" json:"dd_service,omitempty"`
		// ServiceVersion depends on system, can be Git Tag or API version
		//
		// Deprecated: The whole module is deprecated. There is no config struct. Start in
		// github.com/coopnorge/go-telemetry-lib reads the standard OTEL_* environment variables.
		ServiceVersion string `mapstructure:"dd_version" json:"dd_service_version,omitempty"`
		// DSD Socket path for DD StatsD, important to have unix prefix for that value, example: unix:///var/run/dd/dsd.socket
		//
		// Deprecated: The whole module is deprecated. There is no config struct. Start in
		// github.com/coopnorge/go-telemetry-lib reads the standard OTEL_* environment variables.
		DSD string `mapstructure:"dd_dogstatsd_url" json:"dd_dsd,omitempty"`
		// APM Socket path for apm and profiler, unix prefix recommended, but not required, example: unix:///var/run/dd/apm.socket
		//
		// Deprecated: The whole module is deprecated. There is no config struct. Start in
		// github.com/coopnorge/go-telemetry-lib reads the standard OTEL_* environment variables.
		APM string `mapstructure:"dd_trace_agent_url" json:"dd_apm,omitempty"`
		// EnableExtraProfiling flag enables more optional profilers not recommended for production.
		//
		// Deprecated: The whole module is deprecated. There is no config struct. Start in
		// github.com/coopnorge/go-telemetry-lib reads the standard OTEL_* environment variables.
		EnableExtraProfiling bool `mapstructure:"dd_enable_extra_profiling" json:"dd_enable_extra_profiling,omitempty"`
	}
)

// IsDataDogConfigValid method to verify if configuration values are correct
//
// Deprecated: Use Validate().
// The whole module is deprecated. There is no config struct. Start in
// github.com/coopnorge/go-telemetry-lib reads the standard OTEL_* environment variables.
func (d DatadogConfig) IsDataDogConfigValid() bool {
	if err := d.Validate(); err != nil {
		return false
	}

	return true
}

// Validate the DatadogConfig. Returns the first error found, returns nil if
// the configuration is good.
//
// Deprecated: The whole module is deprecated. There is no config struct. Start in
// github.com/coopnorge/go-telemetry-lib reads the standard OTEL_* environment variables.
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
// Deprecated: The whole module is deprecated. There is no config struct. Start in
// github.com/coopnorge/go-telemetry-lib reads the standard OTEL_* environment variables.
func (d DatadogConfig) GetEnv() string {
	return d.Env
}

// GetService how must be service called and displayed in Datadog system
//
// Deprecated: The whole module is deprecated. There is no config struct. Start in
// github.com/coopnorge/go-telemetry-lib reads the standard OTEL_* environment variables.
func (d DatadogConfig) GetService() string {
	return d.Service
}

// GetServiceVersion depends on system, can be Git Tag or API version
//
// Deprecated: The whole module is deprecated. There is no config struct. Start in
// github.com/coopnorge/go-telemetry-lib reads the standard OTEL_* environment variables.
func (d DatadogConfig) GetServiceVersion() string {
	return d.ServiceVersion
}

// GetDsdEndpoint Socket path or URL for DD StatsD
// For unix sockets, the unix-scheme prefix is required.
// Example: unix:///var/run/dd/dsd.socket
//
// Deprecated: The whole module is deprecated. There is no config struct. Start in
// github.com/coopnorge/go-telemetry-lib reads the standard OTEL_* environment variables.
func (d DatadogConfig) GetDsdEndpoint() string {
	return d.DSD
}

// GetApmEndpoint Socket path or URL for APM and profiler
// For unix sockets, the unix-scheme prefix is not needed, but it is recommended to include it.
// Example: unix:///var/run/dd/apm.socket
// Example: http://my-agent:1234
//
// Deprecated: The whole module is deprecated. There is no config struct. Start in
// github.com/coopnorge/go-telemetry-lib reads the standard OTEL_* environment variables.
func (d DatadogConfig) GetApmEndpoint() string {
	return d.APM
}

// IsExtraProfilingEnabled return true if profilers not recommended for production are enabled.
//
// Deprecated: The whole module is deprecated. There is no config struct. Start in
// github.com/coopnorge/go-telemetry-lib reads the standard OTEL_* environment variables.
func (d DatadogConfig) IsExtraProfilingEnabled() bool {
	return d.EnableExtraProfiling
}

// LoadDatadogConfigFromEnvVars loads a new DatadogConfig from known environment-variables.
//
// Deprecated: The whole module is deprecated. There is no config struct. Start in
// github.com/coopnorge/go-telemetry-lib reads the standard OTEL_* environment variables.
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
