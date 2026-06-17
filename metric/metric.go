// Package metric implements custom metrics with Dogstatsd
//
// Deprecated: use metrics package instead
package metric

import (
	"context"
	"fmt"
	"strings"

	"github.com/coopnorge/go-logger"

	"github.com/iancoleman/strcase"
)

type (
	// Name must be in specific format like cart.amount, request.my_request.x etc
	//
	// Deprecated: Use metrics package instead
	Name string
	// Type is an Enum type for metric types
	//
	// Deprecated: Use metrics package instead
	Type byte
	// Tag categories metric value with Name for Value to display category like PaymentID.
	//
	// Deprecated: Use metrics package instead
	Tag struct {
		Name  string
		Value string
	}

	// BaseMetricCollector ...
	//
	// Deprecated: Use metrics package instead
	BaseMetricCollector struct {
		// Deprecated: Use metrics package instead
		DatadogMetrics DatadogMetricsClient
	}

	// Data for metrics
	//
	// Deprecated: Use metrics package instead
	Data struct {
		// Deprecated: Use metrics package instead
		Name Name
		// Deprecated: Use metrics package instead
		Type Type
		// Deprecated: Use metrics package instead
		Value float64
		// MetricTags level empty if no categories required to relate metric
		//
		// Deprecated: Use metrics package instead
		MetricTags []Tag
	}
)

const (
	// MetricTypeCountEvents Datadog will aggregate events to show how many events happened in second
	//
	// Deprecated: Use metrics package instead
	MetricTypeCountEvents Type = iota
	// MetricTypeEvent send single event to Datadog
	//
	// Deprecated: Use metrics package instead
	MetricTypeEvent
	// MetricTypeMeasurement aggregates value of metrics in Datadog for measuring it, like memory or cart value
	//
	// Deprecated: Use metrics package instead
	MetricTypeMeasurement
)

// NewBaseMetricCollector instance
//
// Deprecated: Use coopdatadog.Start() instead.
func NewBaseMetricCollector(dm *DatadogMetrics) *BaseMetricCollector {
	return &BaseMetricCollector{DatadogMetrics: dm}
}

// AddMetric related to name with given value
//
// Deprecated: Use functions from the metrics package
func (m BaseMetricCollector) AddMetric(ctx context.Context, d Data) {
	if m.DatadogMetrics == nil || m.DatadogMetrics.GetClient() == nil {
		return
	}

	var metricTags []string
	for _, t := range d.MetricTags {
		tagName := strings.ToLower(strcase.ToKebab(t.Name))
		metricTags = append(metricTags, fmt.Sprintf("%s:%s", tagName, t.Value))
	}

	metricName := fmt.Sprintf("%s.%s", m.DatadogMetrics.GetServiceNamePrefix(), d.Name)

	var metricCollectionErr error
	switch d.Type {
	case MetricTypeEvent:
		metricCollectionErr = m.DatadogMetrics.GetClient().Incr(metricName, metricTags, 1)
	case MetricTypeMeasurement:
		metricCollectionErr = m.DatadogMetrics.GetClient().Gauge(metricName, d.Value, metricTags, 1)
	case MetricTypeCountEvents:
		metricCollectionErr = m.DatadogMetrics.GetClient().Count(metricName, int64(d.Value), metricTags, 1)
	}

	if metricCollectionErr != nil {
		logger.WithContext(ctx).WithError(metricCollectionErr).Errorf("Failed to collect metrics metricData for Name=%s", metricName)
	}
}

// GracefulShutdown flushes and closes Datadog client
// ensuring that all metrics are sent before the program exits
//
// Deprecated: Use the StopFunc returned coopdatadog.Start() instead.
func (m BaseMetricCollector) GracefulShutdown() {
	if m.DatadogMetrics == nil || m.DatadogMetrics.GetClient() == nil {
		return
	}

	err := m.DatadogMetrics.GetClient().Flush()
	if err != nil {
		logger.Warn("cannot flush Datadog client", err)
	}

	err = m.DatadogMetrics.GetClient().Close()
	if err != nil {
		logger.Warn("cannot close Datadog client", err)
	}
}
