package metrics_test

import (
	"testing"

	"github.com/coopnorge/go-datadog-lib/v2/metrics"
)

func TestUninitlizedMetrics(t *testing.T) {
	t.Parallel()
	// Assert that we can unit-test some code that does not initialize the metrics-package.
	Count("my.metric", 1)
	Gauge("my.gauge", 42.0)
	Incr("my.counter")
	Decr("my.counter")
	Histogram("my.histogram", 100)
	Distribution("my.distribution", 100)
	Set("my.set", "value")
	SimpleEvent("title", "text")
	Count("metric.with.options", 1, WithTags("tag1:value1"))
	Gauge("gauge.with.options", 42.0, WithTags("service", "test"), WithSampleRate(0.5))
}
