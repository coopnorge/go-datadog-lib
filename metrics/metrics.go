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

	clientMu     sync.Mutex
	statsdClient statsd.ClientInterface
	opts         *options
)

func init() {
	// init should initialize the global variables with instances that does not cause panic.
	// These values should only be used when unit-testing code that does not want to call `GlobalSetup`, and not set environment-variables.
	// Any calls to `GlobalSetup` will override this no-op client.
	setNoOpClient()
}

// GlobalSetup configures the Dogstatsd Client. GlobalSetup is intended to be
// called from coopdatadog.Start(), but can be called directly.
func GlobalSetup(options ...Option) error {
	setupOnce.Do(func() {
		if internal.IsDatadogDisabled() {
			setNoOpClient()
			return
		}

		localOpts, err := resolveOptions(options)
		if err != nil {
			setupErr = err
			return
		}

		localClient, err := statsd.New(localOpts.dsdEndpoint, statsd.WithTags(localOpts.tags))
		if setupErr != nil {
			setupErr = err
			return
		}

		setClientInternal(localClient, localOpts)
	})
	return setupErr
}

func setNoOpClient() {
	setClientInternal(&statsd.NoOpClient{}, defaultOptions())
}

func setClientInternal(client statsd.ClientInterface, options *options) {
	clientMu.Lock()
	defer clientMu.Unlock()
	statsdClient = client
	opts = options
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
func Gauge(name string, value float64, tags ...string) {
	err := statsdClient.Gauge(name, value, tags, opts.metricSampleRate)
	if err != nil {
		opts.errorHandler(fmt.Errorf("failed to send Gauge: %w", err))
	}
}

// Count tracks how many times something happened per second.
func Count(name string, value int64, tags ...string) {
	err := statsdClient.Count(name, value, tags, opts.metricSampleRate)
	if err != nil {
		opts.errorHandler(fmt.Errorf("failed to to send Count: %w", err))
	}
}

// Histogram tracks the statistical distribution of a set of values on each host.
func Histogram(name string, value float64, tags ...string) {
	err := statsdClient.Histogram(name, value, tags, opts.metricSampleRate)
	if err != nil {
		opts.errorHandler(fmt.Errorf("failed to to send Histogram: %w", err))
	}
}

// Distribution tracks the statistical distribution of a set of values across your infrastructure.
func Distribution(name string, value float64, tags ...string) {
	err := statsdClient.Distribution(name, value, tags, opts.metricSampleRate)
	if err != nil {
		opts.errorHandler(fmt.Errorf("failed to to send Distribution: %w", err))
	}
}

// Decr is just Count of -1
func Decr(name string, tags ...string) {
	Count(name, -1, tags...)
}

// Incr is just Count of 1
func Incr(name string, tags ...string) {
	Count(name, 1, tags...)
}

// Set counts the number of unique elements in a group.
func Set(name string, value string, tags ...string) {
	err := statsdClient.Set(name, value, tags, opts.metricSampleRate)
	if err != nil {
		opts.errorHandler(fmt.Errorf("failed to to send Set: %w", err))
	}
}

// Timing sends timing information, it is an alias for TimeInMilliseconds
func Timing(name string, value time.Duration, tags ...string) {
	TimeInMilliseconds(name, value.Seconds()*1000, tags...)
}

// TimeInMilliseconds sends timing information in milliseconds.
func TimeInMilliseconds(name string, value float64, tags ...string) {
	err := statsdClient.TimeInMilliseconds(name, value, tags, opts.metricSampleRate)
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
