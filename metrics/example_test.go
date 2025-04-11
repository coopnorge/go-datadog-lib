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

	metricsWithSampleRate, err := metrics.NewClient(metrics.WithTags("a"))
	if err != nil {
		panic(err)
	}

	metrics.Incr("my-metric")
	metricsWithSampleRate.Incr("my-sampled-metric")

	return nil
}
