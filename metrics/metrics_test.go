package metrics_test

import (
	"testing"

	"github.com/coopnorge/go-datadog-lib/v2/metrics"
)

func TestUninitlizedMetrics(t *testing.T) {
	t.Parallel()
	// Assert that we can unit-test some code that does not initialize the metrics-package.
	metrics.Count("my.metric", 1)
}
