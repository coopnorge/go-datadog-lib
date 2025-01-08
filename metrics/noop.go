package metrics

import (
	"time"

	"github.com/DataDog/datadog-go/v5/statsd"
)

// Verify that Client implements the ClientInterface.
var _ statsd.ClientInterface = &noopClient{}

type noopClient struct{}

// Close implements statsd.ClientInterface.
func (n *noopClient) Close() error {
	return nil
}

// Count implements statsd.ClientInterface.
func (n *noopClient) Count(_ string, _ int64, _ []string, _ float64) error {
	return nil
}

// CountWithTimestamp implements statsd.ClientInterface.
func (n *noopClient) CountWithTimestamp(_ string, _ int64, _ []string, _ float64, _ time.Time) error {
	return nil
}

// Decr implements statsd.ClientInterface.
func (n *noopClient) Decr(_ string, _ []string, _ float64) error {
	return nil
}

// Distribution implements statsd.ClientInterface.
func (n *noopClient) Distribution(_ string, _ float64, _ []string, _ float64) error {
	return nil
}

// Event implements statsd.ClientInterface.
func (n *noopClient) Event(_ *statsd.Event) error {
	return nil
}

// Flush implements statsd.ClientInterface.
func (n *noopClient) Flush() error {
	return nil
}

// Gauge implements statsd.ClientInterface.
func (n *noopClient) Gauge(_ string, _ float64, _ []string, _ float64) error {
	return nil
}

// GaugeWithTimestamp implements statsd.ClientInterface.
func (n *noopClient) GaugeWithTimestamp(_ string, _ float64, _ []string, _ float64, _ time.Time) error {
	return nil
}

// GetTelemetry implements statsd.ClientInterface.
func (n *noopClient) GetTelemetry() statsd.Telemetry {
	return statsd.Telemetry{}
}

// Histogram implements statsd.ClientInterface.
func (n *noopClient) Histogram(_ string, _ float64, _ []string, _ float64) error {
	return nil
}

// Incr implements statsd.ClientInterface.
func (n *noopClient) Incr(_ string, _ []string, _ float64) error {
	return nil
}

// IsClosed implements statsd.ClientInterface.
func (n *noopClient) IsClosed() bool {
	return true
}

// ServiceCheck implements statsd.ClientInterface.
func (n *noopClient) ServiceCheck(_ *statsd.ServiceCheck) error {
	return nil
}

// Set implements statsd.ClientInterface.
func (n *noopClient) Set(_ string, _ string, _ []string, _ float64) error {
	return nil
}

// SimpleEvent implements statsd.ClientInterface.
func (n *noopClient) SimpleEvent(_ string, _ string) error {
	return nil
}

// SimpleServiceCheck implements statsd.ClientInterface.
func (n *noopClient) SimpleServiceCheck(_ string, _ statsd.ServiceCheckStatus) error {
	return nil
}

// TimeInMilliseconds implements statsd.ClientInterface.
func (n *noopClient) TimeInMilliseconds(_ string, _ float64, _ []string, _ float64) error {
	return nil
}

// Timing implements statsd.ClientInterface.
func (n *noopClient) Timing(_ string, _ time.Duration, _ []string, _ float64) error {
	return nil
}
