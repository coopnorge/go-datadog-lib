// Package metric implements custom metrics with Dogstatsd
//
// Deprecated: The whole module is deprecated. Use instruments created once from
// otel.Meter in go.opentelemetry.io/otel (package
// go.opentelemetry.io/otel/metric). See the root package documentation for the
// migration guide and ADR.
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
	// Deprecated: The whole module is deprecated. Use instruments created once from
	// otel.Meter in go.opentelemetry.io/otel (package go.opentelemetry.io/otel/metric).
	Name string
	// Type is an Enum type for metric types
	//
	// Deprecated: The whole module is deprecated. Use instruments created once from
	// otel.Meter in go.opentelemetry.io/otel (package go.opentelemetry.io/otel/metric).
	Type byte
	// Tag categories metric value with Name for Value to display category like PaymentID.
	//
	// Deprecated: The whole module is deprecated. Use instruments created once from
	// otel.Meter in go.opentelemetry.io/otel (package go.opentelemetry.io/otel/metric).
	Tag struct {
		Name  string
		Value string
	}

	// BaseMetricCollector ...
	//
	// Deprecated: The whole module is deprecated. Use instruments created once from
	// otel.Meter in go.opentelemetry.io/otel (package go.opentelemetry.io/otel/metric).
	BaseMetricCollector struct {
		DatadogMetrics DatadogMetricsClient
	}

	// Data for metrics
	//
	// Deprecated: The whole module is deprecated. Use instruments created once from
	// otel.Meter in go.opentelemetry.io/otel (package go.opentelemetry.io/otel/metric).
	Data struct {
		Name  Name
		Type  Type
		Value float64
		// MetricTags level empty if no categories required to relate metric
		//
		// Deprecated: The whole module is deprecated. Use instruments created once from
		// otel.Meter in go.opentelemetry.io/otel (package go.opentelemetry.io/otel/metric).
		MetricTags []Tag
	}
)

const (
	// MetricTypeCountEvents Datadog will aggregate events to show how many events happened in second
	//
	// Deprecated: The whole module is deprecated. Use instruments created once from
	// otel.Meter in go.opentelemetry.io/otel (package go.opentelemetry.io/otel/metric).
	MetricTypeCountEvents Type = iota
	// MetricTypeEvent send single event to Datadog
	//
	// Deprecated: The whole module is deprecated. Use instruments created once from
	// otel.Meter in go.opentelemetry.io/otel (package go.opentelemetry.io/otel/metric).
	MetricTypeEvent
	// MetricTypeMeasurement aggregates value of metrics in Datadog for measuring it, like memory or cart value
	//
	// Deprecated: The whole module is deprecated. Use instruments created once from
	// otel.Meter in go.opentelemetry.io/otel (package go.opentelemetry.io/otel/metric).
	MetricTypeMeasurement
)

// NewBaseMetricCollector instance
//
// Deprecated: The whole module is deprecated. Use instruments created once from
// otel.Meter in go.opentelemetry.io/otel (package go.opentelemetry.io/otel/metric).
func NewBaseMetricCollector(dm *DatadogMetrics) *BaseMetricCollector {
	return &BaseMetricCollector{DatadogMetrics: dm}
}

// AddMetric related to name with given value
//
// Deprecated: Use functions from the metrics package.
// The whole module is deprecated; use instruments created once from otel.Meter in
// go.opentelemetry.io/otel (package go.opentelemetry.io/otel/metric).
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
// Deprecated: The whole module is deprecated. Use instruments created once from
// otel.Meter in go.opentelemetry.io/otel (package go.opentelemetry.io/otel/metric).
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
