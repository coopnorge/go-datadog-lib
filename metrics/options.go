package metrics

import (
	"fmt"
	"os"

	"github.com/coopnorge/go-datadog-lib/v2/errors"
	"github.com/coopnorge/go-datadog-lib/v2/internal"
	"github.com/coopnorge/go-logger"
)

const (
	defaultEnableMetrics    = true
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
	options := &options{
		errorHandler: func(err error) {
			logger.WithError(err).Error(err.Error())
		},
		dsdEndpoint:      os.Getenv(internal.DatadogDSDEndpoint),
		metricSampleRate: defaultMetricSampleRate,
		tags: []string{
			fmt.Sprintf("environment:%s", os.Getenv(internal.DatadogEnvironment)),
			fmt.Sprintf("service:%s", os.Getenv(internal.DatadogService)),
			fmt.Sprintf("version:%s", os.Getenv(internal.DatadogVersion)),
		},
	}

	for _, option := range opts {
		err = option(options)
		if err != nil {
			return nil, err
		}
	}

	return options, nil
}

// WithTags sets the tags that are sent with every metric, shorthand for
// statsd.WithTags()
func WithTags(tags ...string) Option {
	return func(options *options) error {
		options.tags = append(options.tags, tags...)
		return nil
	}
}

// WithMetricSampling configures how many metrics are sent to Datadog by
// setting a sampling rate.
//
// For high-volume metrics, sampling helps reduce network traffic while
// maintaining statistical accuracy. Datadog automatically scales up the
// received values to compensate for sampling. Rate must be between 0 and 1,
// where 1.0 sends all metrics and 0.5 sends half.
func WithMetricSampling(rate float64) Option {
	return func(options *options) error {
		options.metricSampleRate = rate
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
