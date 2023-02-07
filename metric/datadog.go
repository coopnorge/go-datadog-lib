package metric

import (
    "fmt"
    "github.com/DataDog/datadog-go/statsd"
    "github.com/coopnorge/go-datadog-lib/config"
    "github.com/iancoleman/strcase"
    "strings"
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
        ServicePrefix      string
        DefaultMetricsTags []string
    }
)

// NewDatadogMetrics instance required to have cfg config.DatadogParameters to get information about service and optional orgPrefix to append into for metric name
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
func (d DatadogMetrics) GetClient() statsd.ClientInterface {
    return d.client
}

// GetDefaultTags that will be used in Datadog metrics
func (d DatadogMetrics) GetDefaultTags() []string {
    return d.DefaultMetricsTags
}

// GetServiceNamePrefix for metric name
func (d DatadogMetrics) GetServiceNamePrefix() string {
    return d.ServicePrefix
}
