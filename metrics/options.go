package metrics

import (
	"fmt"
	"os"

	"github.com/coopnorge/go-datadog-lib/v2/errors"
	"github.com/coopnorge/go-datadog-lib/v2/internal"
	"github.com/coopnorge/go-logger"
)

const (
	defaultMetricSampleRate = 1
)

// Option is used to configure the behaviour of the metrics integration.
type Option func(*options) error

type options struct {
	errorHandler     errors.ErrorHandler
	dsdEndpoint      string
	metricSampleRate float64
	tags             []string
}

// MetricOpts represents a configuration option for metrics.
type MetricOpts func(*metricOpts)

// metricOpts tags and sample rate opts.
type metricOpts struct {
	tags       []string
	sampleRate float64
}

func resolveOptions(opts []Option) (*options, error) {
	err := internal.VerifyEnvVarsSet(
		internal.DatadogDSDEndpoint,
		internal.DatadogEnvironment,
		internal.DatadogService,
		internal.DatadogVersion,
	)
	if err != nil {
		return nil, err
	}

	options := defaultOptions()
	// Apply default options when resolving real options
	options.dsdEndpoint = os.Getenv(internal.DatadogDSDEndpoint)
	options.tags = []string{
		fmt.Sprintf("environment:%s", os.Getenv(internal.DatadogEnvironment)),
		fmt.Sprintf("service:%s", os.Getenv(internal.DatadogService)),
		fmt.Sprintf("version:%s", os.Getenv(internal.DatadogVersion)),
	}

	for _, option := range opts {
		err = option(options)
		if err != nil {
			return nil, err
		}
	}

	return options, nil
}

func defaultOptions() *options {
	return &options{
		errorHandler: func(err error) {
			logger.WithError(err).Error(err.Error())
		},
		metricSampleRate: defaultMetricSampleRate,
	}
}

// parseMetricOptions parses metricOpts
func parseMetricOpts(options ...MetricOpts) metricOpts {
	result := metricOpts{
		tags:       make([]string, 0),
		sampleRate: opts.metricSampleRate,
	}

	for _, opt := range options {
		opt(&result)
	}

	return result
}

// WithGlobalTags sets the tags that are sent with every metric, shorthand for
// statsd.WithTags()
func WithGlobalTags(tags ...string) Option {
	return func(options *options) error {
		options.tags = append(options.tags, tags...)
		return nil
	}
}

// WithErrorHandler allows for setting a custom ErrorHandler to be called on
// function that may error but does not return an error
func WithErrorHandler(handler errors.ErrorHandler) Option {
	return func(options *options) error {
		options.errorHandler = handler
		return nil
	}
}

// WithTags sets the tags that are sent with specific metric
func WithTags(tags ...string) MetricOpts {
	return func(o *metricOpts) {
		for _, tag := range tags {
			if tag != "" { // ignoring empty tags
				o.tags = append(o.tags, tag)
			}
		}
	}
}

// WithSampleRate sets the sample rate
func WithSampleRate(rate float64) MetricOpts {
	return func(o *metricOpts) {
		if rate < 0 {
			rate = 0
		} else if rate > 1 {
			rate = 1
		}
		o.sampleRate = rate
	}
}
