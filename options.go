package coopdatadog

import (
	"github.com/coopnorge/go-datadog-lib/v2/errors"
	"github.com/coopnorge/go-datadog-lib/v2/internal"
	"github.com/coopnorge/go-logger"
)

const (
	defaultEnableExtraProfiling = false
)

// options is the internal configuration for the Datadog integration
type options struct {
	enableExtraProfiling bool
	errorHandler         errors.ErrorHandler
}

func resolveOptions(opts []Option) (*options, error) {
	options := &options{
		enableExtraProfiling: defaultEnableExtraProfiling,
		errorHandler: func(err error) {
			logger.WithError(err).Error(err.Error())
		},
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
type Option func(*options) error

func withConfigFromEnvVars() Option {
	return func(options *options) error {
		options.enableExtraProfiling = internal.GetBool(internal.DatadogEnableExtraProfiling, options.enableExtraProfiling)
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
