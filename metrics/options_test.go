package metrics

import (
	"math"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/stretchr/testify/assert"
)

// TestWithTagValidation tests the WithTag.
func TestWithTagValidation(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		key     string
		value   string
		wantErr bool
	}{
		{"Empty key", "", "value", true},
		{"Valid tag", "key", "value", false},
		{"Invalid key chars", "key:with:colons", "value", true},
		{"URL in value", "key", "https://www.coop.no/some?url=with&different=characters#foobar", false},
		{"Long key", strings.Repeat("a", 250), "value", true},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			opt := WithTag(tc.key, tc.value)
			options := &options{}
			err := opt(options)

			if tc.wantErr {
				require.Error(t, err, "WithTag should return error for invalid input")
			} else {
				require.NoError(t, err, "WithTag should not return error for valid input")
				require.Contains(t, options.tags, tc.key+":"+tc.value, "Tag should be added to options")
			}
		})
	}
}

// TestWithSampleRateValidation tests the WithSampleRate.
func TestWithSampleRateValidation(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		rate    float64
		wantErr bool
	}{
		{"Valid rate", 0.5, false},
		{"Zero rate", 0.0, false},
		{"Full rate", 1.0, false},
		{"Negative rate", -0.1, true},
		{"Rate above maximum", 1.1, true},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			opt := WithSampleRate(tc.rate)
			options := &options{sampleRate: 1.0}
			err := opt(options)

			if tc.wantErr {
				require.Error(t, err, "WithSampleRate should return error for invalid input")
			} else {
				require.NoError(t, err, "WithSampleRate should not return error for valid input")
				if !math.IsNaN(tc.rate) { // Avoid comparing with NaN
					require.Equal(t, tc.rate, options.sampleRate, "Sample rate should be set correctly")
				}
			}
		})
	}
}

func TestMetricOptions(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		options   []Option
		expectErr bool
	}{
		{
			name:      "Valid tag",
			options:   []Option{WithTag("tag1", "value1")},
			expectErr: false,
		},
		{
			name:      "Empty tag key",
			options:   []Option{WithTag("", "value1")},
			expectErr: true,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			opts := options{}
			err := opts.applyOptions(tc.options)

			if tc.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)

				switch tc.name {
				case "Valid tag":
					require.Contains(t, opts.tags, "tag1:value1", "Tag should be added")
				case "Valid sample rate":
					require.Equal(t, 0.5, opts.sampleRate, "Sample rate should be 0.5")
				}
			}
		})
	}
}
