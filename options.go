package coopdatadog

import "github.com/coopnorge/go-datadog-lib/v2/internal"

const (
	defaultEnableTracing        = true
	defaultEnableProfiling      = true
	defaultEnableExtraProfiling = false
)

// config is the internal configuration for the Datadog integration
type config struct {
	enableTracing        bool
	enableProfiling      bool
	enableExtraProfiling bool
}

func resolveConfig(options []Option) (*config, error) {
	cfg := &config{
		enableTracing:        defaultEnableTracing,
		enableProfiling:      defaultEnableProfiling,
		enableExtraProfiling: defaultEnableExtraProfiling,
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
