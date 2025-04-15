package metrics_test

import (
	"testing"

	"github.com/coopnorge/go-datadog-lib/v2/metrics"
)

func TestUninitlizedMetrics(t *testing.T) {
	t.Parallel()
	// Assert that we can unit-test some code that does not initialize the metrics-package.
	metrics.Count("my.metric", 1)
	metrics.Gauge("my.gauge", 42.0)
	metrics.Incr("my.counter")
	metrics.Decr("my.counter")
	metrics.Histogram("my.histogram", 100)
	metrics.Distribution("my.distribution", 100)
	metrics.Set("my.set", "value")
	metrics.SimpleEvent("title", "text")
	metrics.Count("metric.with.options", 1, metrics.WithTag("tag1", "value1"))
	metrics.Gauge("gauge.with.options", 42.0, metrics.WithTag("service", "test"), metrics.WithSampleRate(0.5))
}
