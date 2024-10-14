package v2

import (
	"context"

	"github.com/coopnorge/go-logger"
)

type Gauge struct {
	metric
}

// Record sends a metric to Statsd.
func (g *Gauge) Record(ctx context.Context, value float64) {
	if ctx.Err() != nil {
		logger.Errorf("Context was cancelled, will not send metric: %v", ctx.Err())
	}

	if g.metric.client == nil {
		return
	}

	err := g.client.Gauge(g.name, value, formatTags(g.tags), 1)
	if err != nil {
		logger.WithError(err).Errorf("Failed to Record Gauge: %s", g.name)
	}
}
