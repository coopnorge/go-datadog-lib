package metrics

import (
	"testing"
)

func TestUninitlizedMetrics(t *testing.T) {
	t.Parallel()
	// Assert that we can unit-test some code that does not initialize the metrics-package.
	Count("my.metric", 1)
}
