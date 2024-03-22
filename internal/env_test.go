package internal_test

import (
	"os"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/coopnorge/go-datadog-lib/v2/internal"
)

func TestIsDatadogEnabledFlagNotSet(t *testing.T) {
	t.Setenv(internal.DatadogEnable, "")
	os.Unsetenv(internal.DatadogEnable)

	assert.True(t, internal.IsDatadogEnabled())
}

func TestIsDatadogEnabledFlagEmpty(t *testing.T) {
	t.Setenv(internal.DatadogEnable, "")

	assert.True(t, internal.IsDatadogEnabled())
}

func TestIsDatadogEnabledFlagTrue(t *testing.T) {
	t.Setenv(internal.DatadogEnable, strconv.FormatBool(true))

	assert.True(t, internal.IsDatadogEnabled())
}

func TestIsDatadogEnabledFlagFalse(t *testing.T) {
	t.Setenv(internal.DatadogEnable, strconv.FormatBool(false))

	assert.False(t, internal.IsDatadogEnabled())
}

func TestIsDatadogEnabledNonBoolValue(t *testing.T) {
	t.Setenv(internal.DatadogEnable, "Hello")

	assert.True(t, internal.IsDatadogEnabled())
}
