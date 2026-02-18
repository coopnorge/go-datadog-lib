package internal_test

import (
	"os"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/coopnorge/go-datadog-lib/v2/internal"
)

func TestIsDatadogDisabledFlagNotSet(t *testing.T) {
	t.Setenv(internal.DatadogDisable, "")
	err := os.Unsetenv(internal.DatadogDisable)
	require.NoError(t, err)

	assert.False(t, internal.IsDatadogDisabled())
}

func TestIsDatadogDisabledFlagEmpty(t *testing.T) {
	t.Setenv(internal.DatadogDisable, "")

	assert.False(t, internal.IsDatadogDisabled())
}

func TestIsDatadogDisabledFlagTrue(t *testing.T) {
	t.Setenv(internal.DatadogDisable, strconv.FormatBool(true))

	assert.True(t, internal.IsDatadogDisabled())
}

func TestIsDatadogDisabledFlagFalse(t *testing.T) {
	t.Setenv(internal.DatadogDisable, strconv.FormatBool(false))

	assert.False(t, internal.IsDatadogDisabled())
}

func TestIsDatadogDisabledNonBoolValue(t *testing.T) {
	t.Setenv(internal.DatadogDisable, "Hello")

	assert.False(t, internal.IsDatadogDisabled())
}
