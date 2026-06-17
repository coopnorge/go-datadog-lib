package metric

import (
	"fmt"
	"strings"

	"github.com/coopnorge/go-datadog-lib/v2/config"

	"github.com/DataDog/datadog-go/v5/statsd"
	"github.com/iancoleman/strcase"
)

type (
	// DatadogMetricsClient ...
	//
	// Deprecated: Use metrics package instead
	DatadogMetricsClient interface {
		// GetClient statsd client
		//
		// Deprecated: Use metrics package instead
		GetClient() statsd.ClientInterface
		// GetDefaultTags that will be used in Datadog metrics
		//
		// Deprecated: Use metrics package instead
		GetDefaultTags() []string
		// GetServiceNamePrefix for metric name
		//
		// Deprecated: Use metrics package instead
		GetServiceNamePrefix() string
	}

	// DatadogMetrics ready to use client to send statsd metrics
	//
	// Deprecated: Use metrics package instead
	DatadogMetrics struct {
		client        *statsd.Client
		ServicePrefix string
		// Deprecated: Use metrics package instead
		DefaultMetricsTags []string
	}
)

// NewDatadogMetrics instance required to have cfg config.DatadogParameters to get information about service and optional orgPrefix to append into for metric name
//
// Deprecated: Use coopdatadog.Start() instead.
func NewDatadogMetrics(cfg config.DatadogParameters, orgPrefix string) (*DatadogMetrics, error) {
	var ddClient *statsd.Client
	var ddClientErr error

	ddClient, ddClientErr = statsd.New(cfg.GetDsdEndpoint())
	if ddClientErr != nil {
		return nil, fmt.Errorf("datadog statsd client initialize with socket(%s) - error %v", cfg.GetDsdEndpoint(), ddClientErr)
	}

	dm := &DatadogMetrics{
		client: ddClient,
		ServicePrefix: fmt.Sprintf(
			"%s.%s",
			strings.ToLower(strcase.ToSnake(orgPrefix)),
			strings.ToLower(strcase.ToSnake(cfg.GetService())),
		),
		DefaultMetricsTags: []string{
			fmt.Sprintf("environment:%s", cfg.GetEnv()),
			fmt.Sprintf("service:%s", cfg.GetService()),
			fmt.Sprintf("version:%s", cfg.GetServiceVersion()),
		},
	}

	return dm, nil
}

// GetClient statsd client
//
// Deprecated: Use metrics package instead
func (d DatadogMetrics) GetClient() statsd.ClientInterface {
	return d.client
}

// GetDefaultTags that will be used in Datadog metrics
//
// Deprecated: Use metrics package instead
func (d DatadogMetrics) GetDefaultTags() []string {
	return d.DefaultMetricsTags
}

// GetServiceNamePrefix for metric name
//
// Deprecated: Use metrics package instead
func (d DatadogMetrics) GetServiceNamePrefix() string {
	return d.ServicePrefix
}
