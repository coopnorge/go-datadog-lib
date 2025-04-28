package metrics

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// This test verifies that we can set a global error-handler, but then override it on a per-metric basis.
func TestOverrideGlobalValues(t *testing.T) {
	oldGlobalOpts := globalOpts
	t.Cleanup(func() { globalOpts = oldGlobalOpts })

	optionThatCausesError := Option(func(*options) error { return fmt.Errorf("always return error") })

	globalCount := 0
	globalErrHandler := func(_ error) { globalCount++ }

	localCount := 0
	localErrHandler := func(_ error) { localCount++ }

	globalOpts = &options{errorHandler: globalErrHandler, sampleRate: 1.0}

	Gauge("some_metric", 1.0, optionThatCausesError)
	assert.Equal(t, 1, globalCount)
	assert.Equal(t, 0, localCount)

	Gauge("some_metric", 1.0, optionThatCausesError, WithErrorHandler(localErrHandler))
	assert.Equal(t, 1, globalCount)
	assert.Equal(t, 1, localCount)

	Gauge("some_metric", 1.0, optionThatCausesError)
	assert.Equal(t, 2, globalCount)
	assert.Equal(t, 1, localCount)
}

func TestGetLocalOpts(t *testing.T) {
	oldGlobalOpts := globalOpts
	t.Cleanup(func() { globalOpts = oldGlobalOpts })

	globalOpts = &options{
		errorHandler: func(_ error) {
			_ = "global error handler"
		},
		sampleRate: 0.8,
		tags:       []string{"global tag"},
	}

	localOpts := getLocalOpts()
	require.NotSame(t, globalOpts, localOpts, "globalOpts and localOpts are pointing to the same memory")

	assert.Len(t, globalOpts.tags, 1, "globalOpts no longer have 1 tag")
	assert.Len(t, localOpts.tags, 0, "localOpts copied the global tag, but it should not")

	err := localOpts.applyOptions([]Option{
		WithErrorHandler(func(_ error) {
			_ = "local error handler"
		}),
		WithSampleRate(0.7),
		WithTag("some_key", "some_value"),
	})
	require.NoError(t, err)

	assert.NotSame(t, &globalOpts.errorHandler, &localOpts.errorHandler)
	assert.NotEqual(t, globalOpts.sampleRate, localOpts.sampleRate)
	assert.NotEqual(t, globalOpts.tags, localOpts.tags)
}
