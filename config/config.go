package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
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
		//
		// Deprecated: Use Validate()
		IsDataDogConfigValid() bool
		// Validate the DatadogConfig. Returns the first error found, returns nil if
		// the configuration is good.
		Validate() error
	}

	// DatadogConfig that required to connect to Datadog Agent
	DatadogConfig struct {
		// Env where application is executed, dev, production, staging etc
		Env string `mapstructure:"dd_env" json:"dd_env,omitempty"`
		// Service how must be service called and displayed in Datadog system
		Service string `mapstructure:"dd_service" json:"dd_service,omitempty"`
		// ServiceVersion depends on system, can be Git Tag or API version
		ServiceVersion string `mapstructure:"dd_version" json:"dd_service_version,omitempty"`
		// DSD Socket path for DD StatsD, important to have unix prefix for that value, example: unix:///var/run/dd/dsd.socket
		DSD string `mapstructure:"dd_dogstatsd_url" json:"dd_dsd,omitempty"`
		// APM Socket path for apm and profiler, unix prefix not needed, example: /var/run/dd/apm.socket
		APM string `mapstructure:"dd_trace_agent_url" json:"dd_apm,omitempty"`
		// EnableExtraProfiling flag enables more optional profilers not recommended for production.
		EnableExtraProfiling bool `mapstructure:"dd_enable_extra_profiling" json:"dd_enable_extra_profiling,omitempty"`
	}
)

// IsDataDogConfigValid method to verify if configuration values are correct
//
// Deprecated: Use Validate()
func (d *DatadogConfig) IsDataDogConfigValid() bool {
	if err := d.Validate(); err != nil {
		return false
	}

	return true
}

// Validate the DatadogConfig. Returns the first error found, returns nil if
// the configuration is good.
func (d *DatadogConfig) Validate() error {
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
func (d *DatadogConfig) GetEnv() string {
	return d.Env
}

// GetService how must be service called and displayed in Datadog system
func (d *DatadogConfig) GetService() string {
	return d.Service
}

// GetServiceVersion depends on system, can be Git Tag or API version
func (d *DatadogConfig) GetServiceVersion() string {
	return d.ServiceVersion
}

// GetDsdEndpoint Socket path or URL for DD StatsD
// for socket important to have unix prefix for that value, example: unix:///var/run/dd/dsd.socket
func (d *DatadogConfig) GetDsdEndpoint() string {
	return d.DSD
}

// GetApmEndpoint Socket path or URL for APM and profiler
// unix prefix not needed, example: /var/run/dd/apm.socket
func (d *DatadogConfig) GetApmEndpoint() string {
	return d.APM
}

// IsExtraProfilingEnabled return true if profilers not recommended for production are enabled.
func (d *DatadogConfig) IsExtraProfilingEnabled() bool {
	return d.EnableExtraProfiling
}

// UnmarshalJSON implements encoding/json/decode.go Unmarshaler
// Allows support native json tag parsing of DatadogConfig.
func (d *DatadogConfig) UnmarshalJSON(b []byte) error {
	var tmpCfgMap map[string]json.RawMessage
	if err := json.Unmarshal(b, &tmpCfgMap); err != nil {
		return errors.Join(err, errors.New("failed to unmarshal DatadogConfig"))
	}

	tmpCfgType := reflect.TypeOf(*d)
	tmpCfg := reflect.New(reflect.TypeOf(*d)).Elem()

	for jsonField, jsonValue := range tmpCfgMap {
		for i := 0; i < tmpCfgType.NumField(); i++ {
			if !strings.Contains(fmt.Sprintf("%v", tmpCfgType.Field(i).Tag), jsonField) {
				continue
			}

			// Extract from json field string value that must be parsed as boolean type
			if jsonField == "dd_enable_extra_profiling" {
				var strValue string
				if err := json.Unmarshal(jsonValue, &strValue); err != nil {
					return err
				}

				strBool, strBoolParseErr := strconv.ParseBool(strValue)
				if strBoolParseErr != nil {
					return errors.Join(strBoolParseErr, errors.New("property to parse config field EnableExtraProfiling to bool"))
				}

				tmpCfg.FieldByName("EnableExtraProfiling").SetBool(strBool)
				break
			}

			// Handle rest of values as string and map them
			property := tmpCfg.FieldByName(tmpCfgType.Field(i).Name)
			if property.IsValid() && property.CanSet() {
				propertyValue := reflect.New(property.Type()).Interface()
				if err := json.Unmarshal(jsonValue, propertyValue); err != nil {
					return err
				}
				property.Set(reflect.ValueOf(propertyValue).Elem())
			}
		}
	}

	if tmpCfg.IsValid() {
		*d = tmpCfg.Interface().(DatadogConfig)
	}

	return nil
}
