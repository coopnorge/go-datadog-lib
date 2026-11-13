package coopdatadog

import (
	"time"

	"github.com/coopnorge/go-datadog-lib/v2/errors" //nolint:staticcheck // internal use of the deprecated module
	"github.com/coopnorge/go-datadog-lib/v2/internal"
	"github.com/coopnorge/go-datadog-lib/v2/metrics" //nolint:staticcheck // internal use of the deprecated module
	"github.com/coopnorge/go-logger"
)

const (
	defaultEnableExtraProfiling = false
	defaultStopTimeout          = 10 * time.Second
)

// options is the internal configuration for the Datadog integration
type options struct {
	enableExtraProfiling bool
	errorHandler         errors.ErrorHandler //nolint:staticcheck // internal use of the deprecated module
	stopTimeout          time.Duration
	metricOptions        []metrics.Option //nolint:staticcheck // internal use of the deprecated module
}

func resolveOptions(opts []Option) (*options, error) {
	options := &options{
		enableExtraProfiling: defaultEnableExtraProfiling,
		errorHandler: func(err error) {
			logger.WithError(err).Error(err.Error())
		},
		stopTimeout: defaultStopTimeout,
	}
	opts = append([]Option{withConfigFromEnvVars()}, opts...)

	for _, option := range opts {
		err := option(options)
		if err != nil {
			return nil, err
		}
	}
	return options, nil
}

// Option is used to configure the behaviour of the Datadog integration.
//
// Deprecated: The whole module is deprecated. There is no direct replacement. Start in
// github.com/coopnorge/go-telemetry-lib reads the standard OTEL_* environment variables.
type Option func(*options) error

func withConfigFromEnvVars() Option {
	return func(options *options) error {
		options.enableExtraProfiling = internal.GetBool(internal.DatadogEnableExtraProfiling, options.enableExtraProfiling)
		return nil
	}
}

// WithMetricsOptions allows for passing the options for setting up metrics
//
// Deprecated: The whole module is deprecated. There is no direct replacement. Start in
// github.com/coopnorge/go-telemetry-lib reads the standard OTEL_* environment variables.
func WithMetricsOptions(metricOptions ...metrics.Option) Option {
	return func(options *options) error {
		options.metricOptions = metricOptions
		return nil
	}
}

// WithErrorHandler allows for setting a custom ErrorHandler to be called on
// function that may error but does not return an error
//
// Deprecated: The whole module is deprecated. There is no direct replacement. Start in
// github.com/coopnorge/go-telemetry-lib reads the standard OTEL_* environment variables.
func WithErrorHandler(handler errors.ErrorHandler) Option {
	return func(options *options) error {
		options.errorHandler = handler
		return nil
	}
}

// WithStopTimeout sets the allowed time for graceful shutdown, defaults to 10
// seconds.
//
// Deprecated: The whole module is deprecated. There is no direct replacement. Start in
// github.com/coopnorge/go-telemetry-lib reads the standard OTEL_* environment variables.
func WithStopTimeout(timeout time.Duration) Option {
	return func(options *options) error {
		options.stopTimeout = timeout
		return nil
	}
}

// WithNoStopTimeout waits for graceful shutdown to complete regardless of how
// long.
//
// Deprecated: The whole module is deprecated. There is no direct replacement. Start in
// github.com/coopnorge/go-telemetry-lib reads the standard OTEL_* environment variables.
func WithNoStopTimeout() Option {
	return func(options *options) error {
		options.stopTimeout = 0
		return nil
	}
}
