package metrics

import (
	"fmt"
	"sync"
	"time"

	"github.com/DataDog/datadog-go/v5/statsd"
	"github.com/coopnorge/go-datadog-lib/v2/errors"
	"github.com/coopnorge/go-datadog-lib/v2/internal"
)

var (
	setupOnce    sync.Once
	setupErr     error
	statsdClient statsd.ClientInterface
	errorHandler errors.ErrorHandler
)

// GlobalSetup configures the Dogstatsd Client. GlobalSetup is intended to be
// called from coopdatadog.Start(), but can be called directly.
func GlobalSetup(options ...Option) error {
	setupOnce.Do(func() {
		if internal.IsDatadogDisabled() {
			statsdClient = &noopClient{}
			return
		}

		var cfg *config
		cfg, setupErr = resolveConfig(options)
		if setupErr != nil {
			return
		}

		if !cfg.enableMetrics {
			statsdClient = &noopClient{}
			return
		}

		statsdOptions := append(cfg.statsdOptions, statsd.WithTags(cfg.tags))

		statsdClient, setupErr = statsd.New(cfg.dsdEndpoint, statsdOptions...)
		if setupErr != nil {
			return
		}
	})
	return setupErr
}

// Flush forces a flush of all the queued dogstatsd payloads.
func Flush() error {
	err := statsdClient.Flush()
	if err != nil {
		return fmt.Errorf("failed to flush: %w", err)
	}
	return nil
}

// Gauge measures the value of a metric at a particular time.
func Gauge(name string, value float64, tags []string, rate float64) {
	err := statsdClient.Gauge(name, value, tags, rate)
	if err != nil {
		errorHandler(fmt.Errorf("failed to send Gauge: %w", err))
	}
}

// GaugeWithTimestamp measures the value of a metric at a given time.
// BETA - Please contact our support team for more information to use this feature: https://www.datadoghq.com/support/
// The value will bypass any aggregation on the client side and agent side, this is
// useful when sending points in the past.
//
// Minimum Datadog Agent version: 7.40.0
func GaugeWithTimestamp(name string, value float64, tags []string, rate float64, timestamp time.Time) {
	err := statsdClient.GaugeWithTimestamp(name, value, tags, rate, timestamp)
	if err != nil {
		errorHandler(fmt.Errorf("failed to send GaugeWithTimestamp: %w", err))
	}
}

// Count tracks how many times something happened per second.
func Count(name string, value int64, tags []string, rate float64) {
	err := statsdClient.Count(name, value, tags, rate)
	if err != nil {
		errorHandler(fmt.Errorf("failed to to send Count: %w", err))
	}
}

// CountWithTimestamp tracks how many times something happened at the given second.
// BETA - Please contact our support team for more information to use this feature: https://www.datadoghq.com/support/
// The value will bypass any aggregation on the client side and agent side, this is
// useful when sending points in the past.
//
// Minimum Datadog Agent version: 7.40.0
func CountWithTimestamp(name string, value int64, tags []string, rate float64, timestamp time.Time) {
	err := statsdClient.CountWithTimestamp(name, value, tags, rate, timestamp)
	if err != nil {
		errorHandler(fmt.Errorf("failed to to send CountWithTimestamp: %w", err))
	}
}

// Histogram tracks the statistical distribution of a set of values on each host.
func Histogram(name string, value float64, tags []string, rate float64) {
	err := statsdClient.Histogram(name, value, tags, rate)
	if err != nil {
		errorHandler(fmt.Errorf("failed to to send Histogram: %w", err))
	}
}

// Distribution tracks the statistical distribution of a set of values across your infrastructure.
//
// It is recommended to use `WithMaxBufferedMetricsPerContext` to avoid dropping metrics at high throughput, `rate` can
// also be used to limit the load. Both options can *not* be used together.
func Distribution(name string, value float64, tags []string, rate float64) {
	err := statsdClient.Distribution(name, value, tags, rate)
	if err != nil {
		errorHandler(fmt.Errorf("failed to to send Distribution: %w", err))
	}
}

// Decr is just Count of -1
func Decr(name string, tags []string, rate float64) {
	Count(name, -1, tags, rate)
}

// Incr is just Count of 1
func Incr(name string, tags []string, rate float64) {
	Count(name, 1, tags, rate)
}

// Set counts the number of unique elements in a group.
func Set(name string, value string, tags []string, rate float64) {
	err := statsdClient.Set(name, value, tags, rate)
	if err != nil {
		errorHandler(fmt.Errorf("failed to to send Set: %w", err))
	}
}

// Timing sends timing information, it is an alias for TimeInMilliseconds
func Timing(name string, value time.Duration, tags []string, rate float64) {
	TimeInMilliseconds(name, value.Seconds()*1000, tags, rate)
}

// TimeInMilliseconds sends timing information in milliseconds.
// It is flushed by statsd with percentiles, mean and other info (https://github.com/etsy/statsd/blob/master/docs/metric_types.md#timing)
func TimeInMilliseconds(name string, value float64, tags []string, rate float64) {
	err := statsdClient.TimeInMilliseconds(name, value, tags, rate)
	if err != nil {
		errorHandler(fmt.Errorf("failed to to send TimeInMilliseconds: %w", err))
	}
}

// Event sends the provided Event.
func Event(e *statsd.Event) {
	err := statsdClient.Event(e)
	if err != nil {
		errorHandler(fmt.Errorf("failed to send Event: %w", err))
	}
}

// SimpleEvent sends an event with the provided title and text.
func SimpleEvent(title, text string) {
	Event(statsd.NewEvent(title, text))
}
