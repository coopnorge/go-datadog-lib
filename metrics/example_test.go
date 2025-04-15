package metrics_test

import (
	"context"

	coopdatadog "github.com/coopnorge/go-datadog-lib/v2"
	"github.com/coopnorge/go-datadog-lib/v2/metrics"
)

func Example() {
	err := run()
	if err != nil {
		panic(err)
	}
}

func run() error {
	stop, err := coopdatadog.Start(context.Background())
	if err != nil {
		return err
	}
	defer func() {
		err := stop()
		if err != nil {
			panic(err)
		}
	}()

	metrics.Incr("my-metric")
	metrics.Count("metric.with.options", 1, metrics.WithTag("tag1", "value1"))
	metrics.Gauge("gauge.with.options", 42.0, metrics.WithTag("tag1", "test1"), metrics.WithSampleRate(0.5))

	return nil
}
