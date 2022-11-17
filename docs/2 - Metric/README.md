# Metric - Datadog StatsD

> Datadog supports custom metrics that you
> can utilize depending on the application.
>
> For example, you could use it to track
> value of cart in side e-commerce shop.
>
> Or you could register events for auth attempts.
>
> All depends on the case and what you're looking forward to
> achieve.

## How to use StatsD in Go - Datadog

There is created an abstract client that simply
connects to Datadog StatsD service.

Also, you will already implement simple metric
the collector that you can extend or just use it to
send your events and measurements.

## How to initialize Go- - Datadog

To prepare the configuration you need to look at the `Setup` section.

After that you can create new instance of Datadog client for DD StatsD.

```go
package your_pkg

import (
	"github.com/coopnorge/go-datadog-lib/config"
	"github.com/coopnorge/go-datadog-lib/metric"
)

func MyServiceContainer(ddCfg *config.DatadogConfig) error {

	// After that you will have pure DD StatsD client
	ddClient := metrics.NewDatadogMetrics(ddCfg)
	
	// If you need simple metric collector then create
	ddMetricCollector := metrics.NewBaseMetricCollector(ddClient)
	// ddMetricCollector -> *BaseMetricCollector allows you to send metrics to Datadog
}
```

## Example how to send metrics

When you have `BaseMetricCollector` from pkg `metrics`
you can call create records in Datadog.

```go
package my_metric

import (
	"context"
	
	"github.com/coopnorge/go-datadog-lib/config"
	"github.com/coopnorge/go-datadog-lib/metric"
)

func Example()  {
	ddClient := metrics.NewDatadogMetrics(new(config.DatadogConfig))
	ddMetricCollector := metrics.NewBaseMetricCollector(ddClient)

	tMetricData := metrics.Data{
		Name:  "RuntimeTest",
		Type:  metrics.MetricTypeEvent,
		Value: float64(42),
		MetricTags: []metrics.MetricTag{
			{MetricTagName: "Show", MetricTagValue: "Case"},
		},
	}
	
	ddMetricCollector.AddMetric(context.Background(), tMetricData)
}
```