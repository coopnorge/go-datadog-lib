// Package metrics implements custom metrics with Dogstatsd
package metrics

import (
	"sync"
	"time"
)

var (
	setupOnce sync.Once
	setupErr  error
	// init should initialize the global variables with instances
	// that does not cause panic. These values should only be used
	// when called from unit-testing code that does not want to
	// set environment-variables and call `GlobalSetup`. Any calls
	// to `GlobalSetup` will override this no-op client.
	globalClient Metrics = noOpClient()
)

// Metrics defines the behaviour of metrics clients
type Metrics interface {
	// Flush forces a flush of all the queued dogstatsd payloads.
	Flush() error

	// Gauge measures the value of a metric at a particular time.
	Gauge(name string, value float64, tags ...string)

	// Count tracks how many times something happened per second.
	Count(name string, value int64, tags ...string)

	// Histogram tracks the statistical distribution of a set of values on each host.
	Histogram(name string, value float64, tags ...string)

	// Distribution tracks the statistical distribution of a set of values across your infrastructure.
	Distribution(name string, value float64, tags ...string)

	// Decr is just Count of -1
	Decr(name string, tags ...string)

	// Incr is just Count of 1
	Incr(name string, tags ...string)

	// Set counts the number of unique elements in a group.
	Set(name string, value string, tags ...string)

	// Timing sends timing information, it is an alias for TimeInMilliseconds
	Timing(name string, value time.Duration, tags ...string)

	// TimeInMilliseconds sends timing information in milliseconds.
	TimeInMilliseconds(name string, value float64, tags ...string)

	// SimpleEvent sends an event with the provided title and text.
	SimpleEvent(title, text string)
}

// GlobalSetup configures the Dogstatsd Client. GlobalSetup is intended to be
// called from coopdatadog.Start(), but can be called directly.
func GlobalSetup(options ...Option) error {
	setupOnce.Do(func() {
		client, err := NewClient(options...)
		if err != nil {
			setupErr = err
		}
		globalClient = client
	})
	return setupErr
}

// Flush forces a flush of all the queued dogstatsd payloads.
func Flush() error {
	return globalClient.Flush()
}

// Gauge measures the value of a metric at a particular time.
func Gauge(name string, value float64, tags ...string) {
	globalClient.Gauge(name, value, tags...)
}

// Count tracks how many times something happened per second.
func Count(name string, value int64, tags ...string) {
	globalClient.Count(name, value, tags...)
}

// Histogram tracks the statistical distribution of a set of values on each host.
func Histogram(name string, value float64, tags ...string) {
	globalClient.Histogram(name, value, tags...)
}

// Distribution tracks the statistical distribution of a set of values across your infrastructure.
func Distribution(name string, value float64, tags ...string) {
	globalClient.Distribution(name, value, tags...)
}

// Decr is just Count of -1
func Decr(name string, tags ...string) {
	globalClient.Decr(name, tags...)
}

// Incr is just Count of 1
func Incr(name string, tags ...string) {
	globalClient.Incr(name, tags...)
}

// Set counts the number of unique elements in a group.
func Set(name string, value string, tags ...string) {
	globalClient.Set(name, value, tags...)
}

// Timing sends timing information, it is an alias for TimeInMilliseconds
func Timing(name string, value time.Duration, tags ...string) {
	globalClient.Timing(name, value, tags...)
}

// TimeInMilliseconds sends timing information in milliseconds.
func TimeInMilliseconds(name string, value float64, tags ...string) {
	globalClient.TimeInMilliseconds(name, value, tags...)
}

// SimpleEvent sends an event with the provided title and text.
func SimpleEvent(title, text string) {
	globalClient.SimpleEvent(title, text)
}
