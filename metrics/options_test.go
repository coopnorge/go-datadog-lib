package metrics

import (
	"math"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/stretchr/testify/assert"
)

func TestWithTags(t *testing.T) {
	options := &options{
		tags: []string{"a", "b"},
	}

	err := WithTags("c")(options)
	assert.NoError(t, err)

	assert.Equal(t, []string{"a", "b", "c"}, options.tags)

	err = WithTags("d")(options)
	assert.NoError(t, err)
	assert.Equal(t, []string{"a", "b", "c", "d"}, options.tags)
}

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
		{"Invalid value chars", "key", "value:with:colons", true},
		{"Long key", strings.Repeat("a", 250), "value", true},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			opt := WithTag(tc.key, tc.value)
			options := &metricOptions{tags: []string{}}
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
			options := &metricOptions{sampleRate: 1.0}
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
		options   []MetricOptions
		expectErr bool
	}{
		{
			name:      "Valid tag",
			options:   []MetricOptions{WithTag("tag1", "value1")},
			expectErr: false,
		},
		{
			name:      "Empty tag key",
			options:   []MetricOptions{WithTag("", "value1")},
			expectErr: true,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			opts := parseMetricOptions(tc.options...)

			if tc.expectErr {
				require.NotNil(t, opts.errorHandler, "Expected error handler to be set for invalid options")

				if opts.errorHandler == nil {
					assert.False(t, true, "Expected error handler to not be set for invalid options")
				}
			} else {
				require.Nil(t, opts.errorHandler, "Expected no error handler for valid options")

				if tc.name == "Valid tag" {
					require.Contains(t, opts.tags, "tag1:value1", "Tag should be added")
				} else if tc.name == "Valid sample rate" {
					require.Equal(t, 0.5, opts.sampleRate, "Sample rate should be 0.5")
				}
			}
		})
	}
}
