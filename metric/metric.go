package metric

import (
    "context"
    "fmt"
    "strings"

    "github.com/coopnorge/go-datadog-lib/tracing"
    "github.com/coopnorge/go-logger"

    "github.com/iancoleman/strcase"
)

type (
    // MetricName must be in specific format like cart.amount, request.my_request.x etc
    MetricName string
    MetricType byte
    // MetricTag categories metric value with MetricTagName for MetricTagValue to display category like PaymentID.
    MetricTag struct {
        MetricTagName  string
        MetricTagValue string
    }

    // BaseMetricCollector ...
    BaseMetricCollector struct {
        DatadogMetrics DatadogMetricsClient
    }

    // Data for metrics
    Data struct {
        Name  MetricName
        Type  MetricType
        Value float64
        // MetricTags level empty if no categories required to relate metric
        MetricTags []MetricTag
    }
)

const (
    // MetricTypeCountEvents Datadog will aggregate events to show how many events happened in second
    MetricTypeCountEvents MetricType = iota
    // MetricTypeEvent send single event to Datadog
    MetricTypeEvent
    // MetricTypeMeasurement aggregates value of metrics in Datadog for measuring it, like memory or cart value
    MetricTypeMeasurement
)

// NewBaseMetricCollector instance
func NewBaseMetricCollector(dm *DatadogMetrics) *BaseMetricCollector {
    return &BaseMetricCollector{DatadogMetrics: dm}
}

// AddMetric related to name with given value
func (m BaseMetricCollector) AddMetric(ctx context.Context, d Data) {
    if m.DatadogMetrics == nil || m.DatadogMetrics.GetClient() == nil {
        return
    }

    var metricTags []string
    for _, t := range d.MetricTags {
        tagName := strings.ToLower(strcase.ToKebab(t.MetricTagName))
        metricTags = append(metricTags, fmt.Sprintf("%s:%s", tagName, t.MetricTagValue))
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
        tracing.LogWithTrace(
            ctx,
            logger.LevelError,
            fmt.Sprintf("Failed to collect metrics metricData for MetricName=%s - error: %v", metricName, metricCollectionErr),
        )
    }
}
