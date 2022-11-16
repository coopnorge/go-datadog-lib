package metrics

import (
	"fmt"
	"strings"

	"github.com/coopnorge/go-datadog-lib/config"
	"github.com/coopnorge/go-logger"

	"github.com/DataDog/datadog-go/statsd"
	"github.com/iancoleman/strcase"
)

type (
	DatadogMetricsClient interface {
		// GetClient statsd client
		GetClient() statsd.ClientInterface
		// GetDefaultTags that will be used in Datadog metrics
		GetDefaultTags() []string
		// GetServiceNamePrefix for metric name
		GetServiceNamePrefix() string
	}

	// DatadogMetrics ready to use client to send statsd metrics
	DatadogMetrics struct {
		client             *statsd.Client
		servicePrefix      string
		defaultMetricsTags []string
	}
)

// NewDatadogMetrics instance
func NewDatadogMetrics(cfg *config.DatadogConfig) *DatadogMetrics {
	var ddClient *statsd.Client
	var ddClientErr error

	ddClient, ddClientErr = statsd.New(cfg.DSD)
	if ddClientErr != nil {
		logger.Errorf("datadog statsd client initialize with socket(%s) - error %w", cfg.DSD, ddClientErr)
		if ddClient, ddClientErr = statsd.New(""); ddClientErr != nil {
			logger.Errorf("datadog statsd self-resolving client initialize - error %w", ddClientErr)
		}
	}

	return &DatadogMetrics{
		client:        ddClient,
		servicePrefix: strings.ToLower(strcase.ToSnake(cfg.Service)),
		defaultMetricsTags: []string{
			fmt.Sprintf("environment:%s", cfg.Env),
			fmt.Sprintf("service:%s", cfg.Service),
			fmt.Sprintf("version:%s", cfg.ServiceVersion),
		},
	}
}

// GetClient statsd client
func (d DatadogMetrics) GetClient() statsd.ClientInterface {
	return d.client
}

// GetDefaultTags that will be used in Datadog metrics
func (d DatadogMetrics) GetDefaultTags() []string {
	return d.defaultMetricsTags
}

// GetServiceNamePrefix for metric name
func (d DatadogMetrics) GetServiceNamePrefix() string {
	return d.servicePrefix
}
