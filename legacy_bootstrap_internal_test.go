package coopdatadog

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNormalizeLegacySocketPath(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"unix:///var/run/datadog/apm.socket", "/var/run/datadog/apm.socket"},
		{"/var/run/datadog/apm.socket", "/var/run/datadog/apm.socket"},
	}
	for _, tc := range tests {
		t.Run(tc.input, func(t *testing.T) {
			got := normalizeLegacySocketPath(tc.input)
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestNormalizeLegacyHTTPAddr(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"http://my-dd-agent:3678", "my-dd-agent:3678"},
		{"http://my-dd-agent", "my-dd-agent"},
		{"http://10.0.0.6", "10.0.0.6"},
	}
	for _, tc := range tests {
		t.Run(tc.input, func(t *testing.T) {
			got := normalizeLegacyHTTPAddr(tc.input)
			assert.Equal(t, tc.want, got)
		})
	}
}
