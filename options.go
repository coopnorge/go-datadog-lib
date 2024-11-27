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

func defaultConfig() *config {
	return &config{
		enableTracing:        defaultEnableTracing,
		enableProfiling:      defaultEnableProfiling,
		enableExtraProfiling: defaultEnableExtraProfiling,
	}
}

// Option is used to configure the behaviour of the Datadog integration.
type Option func(*config)

func withConfigFromEnvVars() Option {
	return func(cfg *config) {
		cfg.enableTracing = getBoolEnv(internal.DatadogEnableTracing, cfg.enableTracing)
		cfg.enableProfiling = getBoolEnv(internal.DatadogEnableProfiling, cfg.enableProfiling)
		cfg.enableExtraProfiling = getBoolEnv(internal.DatadogEnableExtraProfiling, cfg.enableExtraProfiling)
	}
}
