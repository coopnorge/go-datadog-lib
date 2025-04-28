package metrics

import (
	"errors"
	"fmt"
	"os"
	"slices"
	"strings"

	ddErrors "github.com/coopnorge/go-datadog-lib/v2/errors"
	"github.com/coopnorge/go-datadog-lib/v2/internal"
	"github.com/coopnorge/go-logger"
)

const (
	defaultMetricSampleRate = 1
)

// Option is used to configure the behaviour of the metrics integration.
type Option func(*options) error

type options struct {
	dsdEndpoint  string
	errorHandler ddErrors.ErrorHandler
	sampleRate   float64
	tags         []string
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
	options.dsdEndpoint = os.Getenv(internal.DatadogDSDEndpoint)
	// Apply default options when resolving real options
	options.tags = []string{
		fmt.Sprintf("environment:%s", os.Getenv(internal.DatadogEnvironment)),
		fmt.Sprintf("service:%s", os.Getenv(internal.DatadogService)),
		fmt.Sprintf("version:%s", os.Getenv(internal.DatadogVersion)),
	}

	if err := options.applyOptions(opts); err != nil {
		return nil, err
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

// applyOptions applies every option, and returns a combined error of all (if any) errors.
func (opts *options) applyOptions(options []Option) error {
	errs := make([]error, 0, len(options))
	for _, option := range options {
		err := option(opts)
		if err != nil {
			errs = append(errs, err)
		}
	}
	return errors.Join(errs...)
}

// WithErrorHandler allows for setting a custom ErrorHandler to be called on
// function that may error but does not return an error
func WithErrorHandler(handler ddErrors.ErrorHandler) Option {
	return func(options *options) error {
		options.errorHandler = handler
		return nil
	}
}

// WithTag sets a tag that will be sent with a specific metric.
// Parameters:
//   - k: The tag key
//   - v: The tag value
func WithTag(k, v string) Option {
	return func(options *options) error {
		if k == "" {
			return fmt.Errorf("tag key cannot be empty")
		}
		if v == "" {
			return fmt.Errorf("tag value cannot be empty")
		}
		if strings.ContainsAny(k, ":,|=") {
			return fmt.Errorf("tag key contains invalid characters: %s", k)
		}
		if len(k)+len(v)+1 > 200 {
			return fmt.Errorf("tag %s:%s exceeds maximum length", k, v)
		}
		if slices.Contains([]string{"environment", "service", "version"}, strings.ToLower(k)) {
			return fmt.Errorf("tag key '%s' is reserved", k)
		}

		options.tags = append(options.tags, fmt.Sprintf(k+":"+v))
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
func WithSampleRate(rate float64) Option {
	return func(o *options) error {
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
