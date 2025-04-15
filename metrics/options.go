package metrics

import (
	"fmt"
	"os"
	"strings"

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
	sampleRate   float64
	tags         []string
}

// MetricOptions represents a configuration option for metrics.
type MetricOptions func(*metricOptions) error

// metricOptions tags and sample rate opts.
type metricOptions struct {
	errorHandler errors.ErrorHandler
	tags         []string
	sampleRate   float64
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
		sampleRate: defaultMetricSampleRate,
	}
}

// parseMetricOptions processes the provided metric options and returns a configured metricOptions object.
// This function combines globally defined defaults with any user-provided option overrides.
// Notes:
//   - If no options are provided, the returned metricOptions will use:
//   - No tag
//   - Global sample rate (opts.sampleRate) as the default sampling rate
//   - Any provided options will override these defaults for the specific metric
func parseMetricOptions(options ...MetricOptions) metricOptions {
	result := metricOptions{
		tags:       []string{},
		sampleRate: opts.sampleRate, // Use global sample as default
	}

	for _, opt := range options {
		err := opt(&result)
		if err != nil {
			return metricOptions{
				errorHandler: func(err error) {
					logger.WithError(err).Error(err.Error())
				},
			}
		}
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
//   - k: The tag key
//   - v: The tag value
func WithTag(k, v string) MetricOptions {
	return func(o *metricOptions) error {
		if k == "" {
			return fmt.Errorf("tag key cannot be empty")
		}

		if strings.ContainsAny(k, ":,|=") {
			return fmt.Errorf("tag key contains invalid characters: %s", k)
		}

		if strings.ContainsAny(v, ":,|=") {
			return fmt.Errorf("tag value contains invalid characters: %s", v)
		}

		if len(k)+len(v)+1 > 200 {
			return fmt.Errorf("tag %s:%s exceeds maximum length", k, v)
		}

		for _, reservedTag := range []string{"environment", "service", "version"} {
			if strings.ToLower(k) == reservedTag {
				return fmt.Errorf("tag key '%s' is reserved", k)
			}
		}
		o.tags = append(o.tags, fmt.Sprintf(k+":"+v))
		return nil
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
// Returns an error if the rate is invalid (negative)
func WithSampleRate(rate float64) MetricOptions {
	return func(o *metricOptions) error {
		if rate < 0 {
			return fmt.Errorf("sample rate cannot be negative: %f", rate)
		}

		if rate > 1 {
			o.sampleRate = 1.0
			return fmt.Errorf("sample rate %f exceeds maximum of 1.0, capped at 1.0", rate)
		}
		o.sampleRate = rate
		return nil
	}
}
