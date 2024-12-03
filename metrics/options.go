package metrics

import (
	"fmt"
	"os"

	"github.com/coopnorge/go-datadog-lib/v2/errors"
	"github.com/coopnorge/go-datadog-lib/v2/internal"
	"github.com/coopnorge/go-logger"
)

const (
	defaultEnableMetrics = true
)

// Option is used to configure the behaviour of the metrics integration.
type Option func(*config) error

type config struct {
	enableMetrics bool
	errorHandler  errors.ErrorHandler
	dsdEndpoint   string
	tags          []string
}

func resolveConfig(options []Option) (*config, error) {
	err := internal.VerifyEnvVarsSet(
		internal.DatadogDSDEndpoint,
		internal.DatadogEnvironment,
		internal.DatadogService,
		internal.DatadogVersion,
	)
	if err != nil {
		return nil, err
	}
	cfg := &config{
		enableMetrics: internal.GetBool(internal.DatadogEnableMetrics, defaultEnableMetrics),
		errorHandler: func(err error) {
			logger.WithError(err).Error(err.Error())
		},
		dsdEndpoint: os.Getenv(internal.DatadogDSDEndpoint),
		tags: []string{
			fmt.Sprintf("environment:%s", os.Getenv(internal.DatadogEnvironment)),
			fmt.Sprintf("service:%s", os.Getenv(internal.DatadogService)),
			fmt.Sprintf("version:%s", os.Getenv(internal.DatadogVersion)),
		},
	}

	for _, option := range options {
		err = option(cfg)
		if err != nil {
			return nil, err
		}
	}

	return cfg, nil
}

// WithTags sets the tags that are sent with every metric, shorthand for
// statsd.WithTags()
func WithTags(tags ...string) Option {
	return func(cfg *config) error {
		cfg.tags = tags
		return nil
	}
}

// WithErrorHandler allows for setting a custom ErrorHandler to be called on
// function that may error but does not return an error
func WithErrorHandler(handler errors.ErrorHandler) Option {
	return func(cfg *config) error {
		cfg.errorHandler = handler
		return nil
	}
}
