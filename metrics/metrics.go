// Package metrics implements custom metrics with Dogstatsd
package metrics

import (
	"fmt"
	"sync"
	"time"

	"github.com/DataDog/datadog-go/v5/statsd"
	"github.com/coopnorge/go-datadog-lib/v2/internal"
)

var (
	setupOnce sync.Once
	setupErr  error

	// init should initialize the global variables with instances that does not cause panic.
	// These values should only be used when called from unit-testing code that does not want to set environment-variables and call `GlobalSetup`.
	// Any calls to `GlobalSetup` will override this no-op client.
	statsdClient statsd.ClientInterface = &statsd.NoOpClient{}
	opts                                = defaultOptions()
)

// GlobalSetup configures the Dogstatsd Client. GlobalSetup is intended to be
// called from coopdatadog.Start(), but can be called directly.
func GlobalSetup(options ...Option) error {
	setupOnce.Do(func() {
		if internal.IsDatadogDisabled() {
			// Use no-op client initialized by default.
			return
		}

		opts, setupErr = resolveOptions(options)
		if setupErr != nil {
			return
		}

		statsdClient, setupErr = statsd.New(opts.dsdEndpoint, statsd.WithTags(opts.tags))
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
func Gauge(name string, value float64, options ...MetricOpts) {
	metricOpts := parseMetricOpts(options...)
	err := statsdClient.Gauge(name, value, metricOpts.tags, metricOpts.sampleRate)
	if err != nil {
		opts.errorHandler(fmt.Errorf("failed to send Gauge: %w", err))
	}
}

// Count tracks how many times something happened per second.
func Count(name string, value int64, options ...MetricOpts) {
	metricOpts := parseMetricOpts(options...)
	err := statsdClient.Count(name, value, metricOpts.tags, opts.SampleRate)
	if err != nil {
		opts.errorHandler(fmt.Errorf("failed to to send Count: %w", err))
	}
}

// Histogram tracks the statistical distribution of a set of values on each host.
func Histogram(name string, value float64, options ...MetricOpts) {
	metricOpts := parseMetricOpts(options...)
	err := statsdClient.Histogram(name, value, metricOpts.tags, opts.SampleRate)
	if err != nil {
		opts.errorHandler(fmt.Errorf("failed to to send Histogram: %w", err))
	}
}

// Distribution tracks the statistical distribution of a set of values across your infrastructure.
func Distribution(name string, value float64, options ...MetricOpts) {
	metricOpts := parseMetricOpts(options...)
	err := statsdClient.Distribution(name, value, metricOpts.tags, opts.SampleRate)
	if err != nil {
		opts.errorHandler(fmt.Errorf("failed to to send Distribution: %w", err))
	}
}

// Decr is just Count of -1
func Decr(name string, options ...MetricOpts) {
	Count(name, -1, options...)
}

// Incr is just Count of 1
func Incr(name string, options ...MetricOpts) {
	Count(name, 1, options...)
}

// Set counts the number of unique elements in a group.
func Set(name string, value string, options ...MetricOpts) {
	metricOpts := parseMetricOpts(options...)
	err := statsdClient.Set(name, value, metricOpts.tags, opts.SampleRate)
	if err != nil {
		opts.errorHandler(fmt.Errorf("failed to to send Set: %w", err))
	}
}

// Timing sends timing information, it is an alias for TimeInMilliseconds
func Timing(name string, value time.Duration, options ...MetricOpts) {
	TimeInMilliseconds(name, value.Seconds()*1000, options...)
}

// TimeInMilliseconds sends timing information in milliseconds.
func TimeInMilliseconds(name string, value float64, options ...MetricOpts) {
	metricOpts := parseMetricOpts(options...)
	err := statsdClient.TimeInMilliseconds(name, value, metricOpts.tags, opts.SampleRate)
	if err != nil {
		opts.errorHandler(fmt.Errorf("failed to to send TimeInMilliseconds: %w", err))
	}
}

// SimpleEvent sends an event with the provided title and text.
func SimpleEvent(title, text string) {
	err := statsdClient.SimpleEvent(title, text)
	if err != nil {
		opts.errorHandler(fmt.Errorf("failed to send Event: %w", err))
	}
}
