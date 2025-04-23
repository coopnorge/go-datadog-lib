package metrics

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
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
