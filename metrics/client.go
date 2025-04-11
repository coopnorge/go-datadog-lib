package metrics

import (
	"fmt"
	"time"

	"github.com/DataDog/datadog-go/v5/statsd"
	"github.com/coopnorge/go-datadog-lib/v2/internal"
)

var _ Metrics = (*Client)(nil)

// Client implements the Metrics interface
type Client struct {
	statsdClient statsd.ClientInterface
	opts         *options
}

// NewClient creates a metrics client
func NewClient(options ...Option) (*Client, error) {
	if internal.IsDatadogDisabled() {
		// Use no-op client initialized by default.
		return noOpClient(), nil
	}

	opts, err := resolveOptions(options)
	if err != nil {
		return nil, err
	}

	statsdClient, err := statsd.New(opts.dsdEndpoint, statsd.WithTags(opts.tags))
	if err != nil {
		return nil, err
	}

	return &Client{
		statsdClient: statsdClient,
		opts:         opts,
	}, nil
}

// Flush forces a flush of all the queued dogstatsd payloads.
func (c *Client) Flush() error {
	err := c.statsdClient.Flush()
	if err != nil {
		return fmt.Errorf("failed to flush: %w", err)
	}
	return nil
}

// Gauge measures the value of a metric at a particular time.
func (c *Client) Gauge(name string, value float64, tags ...string) {
	err := c.statsdClient.Gauge(name, value, tags, c.opts.metricSampleRate)
	if err != nil {
		c.opts.errorHandler(fmt.Errorf("failed to send Gauge: %w", err))
	}
}

// Count tracks how many times something happened per second.
func (c *Client) Count(name string, value int64, tags ...string) {
	err := c.statsdClient.Count(name, value, tags, c.opts.metricSampleRate)
	if err != nil {
		c.opts.errorHandler(fmt.Errorf("failed to to send Count: %w", err))
	}
}

// Histogram tracks the statistical distribution of a set of values on each host.
func (c *Client) Histogram(name string, value float64, tags ...string) {
	err := c.statsdClient.Histogram(name, value, tags, c.opts.metricSampleRate)
	if err != nil {
		c.opts.errorHandler(fmt.Errorf("failed to to send Histogram: %w", err))
	}
}

// Distribution tracks the statistical distribution of a set of values across your infrastructure.
func (c *Client) Distribution(name string, value float64, tags ...string) {
	err := c.statsdClient.Distribution(name, value, tags, c.opts.metricSampleRate)
	if err != nil {
		c.opts.errorHandler(fmt.Errorf("failed to to send Distribution: %w", err))
	}
}

// Decr is just Count of -1
func (c *Client) Decr(name string, tags ...string) {
	c.Count(name, -1, tags...)
}

// Incr is just Count of 1
func (c *Client) Incr(name string, tags ...string) {
	c.Count(name, 1, tags...)
}

// Set counts the number of unique elements in a group.
func (c *Client) Set(name string, value string, tags ...string) {
	err := c.statsdClient.Set(name, value, tags, c.opts.metricSampleRate)
	if err != nil {
		c.opts.errorHandler(fmt.Errorf("failed to to send Set: %w", err))
	}
}

// Timing sends timing information, it is an alias for TimeInMilliseconds
func (c *Client) Timing(name string, value time.Duration, tags ...string) {
	c.TimeInMilliseconds(name, value.Seconds()*1000, tags...)
}

// TimeInMilliseconds sends timing information in milliseconds.
func (c *Client) TimeInMilliseconds(name string, value float64, tags ...string) {
	err := c.statsdClient.TimeInMilliseconds(name, value, tags, c.opts.metricSampleRate)
	if err != nil {
		c.opts.errorHandler(fmt.Errorf("failed to to send TimeInMilliseconds: %w", err))
	}
}

// SimpleEvent sends an event with the provided title and text.
func (c *Client) SimpleEvent(title, text string) {
	err := c.statsdClient.SimpleEvent(title, text)
	if err != nil {
		c.opts.errorHandler(fmt.Errorf("failed to send Event: %w", err))
	}
}

func noOpClient() *Client {
	return &Client{
		statsdClient: &statsd.NoOpClient{},
		opts:         defaultOptions(),
	}
}
