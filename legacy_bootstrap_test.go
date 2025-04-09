package coopdatadog_test

import (
	"os"
	"testing"

	coopdatadog "github.com/coopnorge/go-datadog-lib/v2"
	"github.com/coopnorge/go-datadog-lib/v2/config"
	"github.com/coopnorge/go-datadog-lib/v2/internal"
	"github.com/stretchr/testify/assert"
)

func TestStartDatadog(t *testing.T) {
	ddCfg := config.DatadogConfig{
		Env:                  "local",
		Service:              "Test-Go-Datadog-lib",
		ServiceVersion:       "na",
		DSD:                  "unix:///tmp/",
		APM:                  "/tmp",
		EnableExtraProfiling: true,
	}

	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		cfg            config.DatadogParameters
		connectionType coopdatadog.ConnectionType
		envVars        func()
		wantErr        bool
	}{
		{
			name:           "Empty config",
			cfg:            config.DatadogConfig{},
			connectionType: coopdatadog.ConnectionTypeHTTP,
			envVars:        func() {},
			wantErr:        true,
		},
		{
			name:           "Socket no env",
			cfg:            ddCfg,
			connectionType: coopdatadog.ConnectionTypeSocket,
			envVars:        func() {},
			wantErr:        false,
		},
		{
			name:           "Socket with env",
			cfg:            ddCfg,
			connectionType: coopdatadog.ConnectionTypeSocket,
			envVars: func() {
				t.Setenv(internal.DatadogAPMEndpoint, "foobar")
			},
			wantErr: false,
		},
		{
			name:           "HTTP no env",
			cfg:            ddCfg,
			connectionType: coopdatadog.ConnectionTypeHTTP,
			envVars:        func() {},
			wantErr:        false,
		},
		{
			name:           "HTTP with env",
			cfg:            ddCfg,
			connectionType: coopdatadog.ConnectionTypeHTTP,
			envVars: func() {
				t.Setenv(internal.DatadogAPMEndpoint, "foobar")
			},
			wantErr: false,
		},
		{
			name: "Auto no env",
			cfg: config.DatadogConfig{
				Env:                  "local",
				Service:              "Test-Go-Datadog-lib",
				ServiceVersion:       "na",
				DSD:                  "unix:///tmp/",
				EnableExtraProfiling: true,
			},
			connectionType: coopdatadog.ConnectionTypeAuto,
			envVars:        func() {},
			wantErr:        true,
		},
		{
			name:           "Auto with env",
			cfg:            ddCfg,
			connectionType: coopdatadog.ConnectionTypeAuto,
			envVars: func() {
				t.Setenv(internal.DatadogAPMEndpoint, "foobar")
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Cleanup(resetEnvVars)

			tt.envVars()
			gotErr := coopdatadog.StartDatadog(tt.cfg, tt.connectionType)
			if tt.wantErr {
				assert.Error(t, gotErr)
			} else {
				assert.NoError(t, gotErr)
			}

			coopdatadog.GracefulDatadogShutdown()
		})
	}
}

func resetEnvVars() {
	os.Unsetenv(internal.DatadogEnvironment)
	os.Unsetenv(internal.DatadogService)
	os.Unsetenv(internal.DatadogVersion)
	os.Unsetenv(internal.DatadogDSDEndpoint)
	os.Unsetenv(internal.DatadogAPMEndpoint)
}
