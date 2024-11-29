package coopdatadog

import (
	"github.com/coopnorge/go-datadog-lib/v2/internal"
	"github.com/coopnorge/go-logger"
)

const (
	defaultEnableTracing        = true
	defaultEnableProfiling      = true
	defaultEnableExtraProfiling = false
)

// ErrorHandler allows for handling of error that cannot be returned to the
// caller
type ErrorHandler func(error)

// config is the internal configuration for the Datadog integration
type config struct {
	enableTracing        bool
	enableProfiling      bool
	enableExtraProfiling bool
	errorHandler         ErrorHandler
}

func resolveConfig(options []Option) (*config, error) {
	cfg := &config{
		enableTracing:        defaultEnableTracing,
		enableProfiling:      defaultEnableProfiling,
		enableExtraProfiling: defaultEnableExtraProfiling,
		errorHandler: func(err error) {
			logger.WithError(err).Error(err.Error())
		},
	}
	options = append([]Option{withConfigFromEnvVars()}, options...)

	for _, option := range options {
		err := option(cfg)
		if err != nil {
			return nil, err
		}
	}
	return cfg, nil
}

// Option is used to configure the behaviour of the Datadog integration.
type Option func(*config) error

func withConfigFromEnvVars() Option {
	return func(cfg *config) error {
		cfg.enableTracing = internal.GetBool(internal.DatadogEnableTracing, cfg.enableTracing)
		cfg.enableProfiling = internal.GetBool(internal.DatadogEnableProfiling, cfg.enableProfiling)
		cfg.enableExtraProfiling = internal.GetBool(internal.DatadogEnableExtraProfiling, cfg.enableExtraProfiling)
		return nil
	}
}

// WithErrorHandler allows for setting a custom ErrorHandler to be called on
// function that may error but does not return an error
func WithErrorHandler(handler ErrorHandler) Option {
	return func(cfg *config) error {
		cfg.errorHandler = handler
		return nil
	}
}
