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
	errorHandler errors.ErrorHandler
	dsdEndpoint  string
	SampleRate   float64
	tags         []string
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
		SampleRate: defaultMetricSampleRate,
	}
}

// parseMetricOpts processes the provided metric options and returns a configured metricOpts object.
// This function combines globally defined defaults with any user-provided option overrides.
// Notes:
//   - If no options are provided, the returned metricOpts will use:
//   - Global tags (opts.tags) as the default tags
//   - Global sample rate (opts.SampleRate) as the default sampling rate
//   - Any provided options will override these defaults for the specific metric
func parseMetricOpts(options ...MetricOpts) metricOpts {
	result := metricOpts{
		tags:       append([]string{}, opts.tags...), // Copy global tags as default
		sampleRate: opts.SampleRate,                  // Use global sample as default
	}

	for _, opt := range options {
		opt(&result)
	}

	return result
}

// WithTags sets the tags that are sent with every metric, shorthand for
// statsd.WithTags()
func WithTags(tags ...string) Option {
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

// WithTag sets a tag that will be sent with a specific metric.
// Parameters:
//   - k: The tag key (if empty, the tag will be ignored)
//   - v: The tag value
func WithTag(k, v string) MetricOpts {
	return func(o *metricOpts) {
		if k == "" {
			return
		}
		o.tags = append(o.tags, fmt.Sprintf(k+":"+v))
	}
}

// WithSampleRate sets the sample rate for metrics collection.
// The sample rate controls what percentage of metrics are actually sent to the backend
// Parameters:
//   - rate: A float between 0 and 1 representing the sampling percentage:
//   - 0: No metrics will be sent (0%)
//   - 1: All metrics will be sent (100%)
//   - between 0 and 1 means that % of metrics will be sent (0.25 = 25%)
//
// Notes:
//   - Values below 0 will be treated as 0 (no metrics sent)
//   - Values above 1 will be treated as 1 (all metrics sent)
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
