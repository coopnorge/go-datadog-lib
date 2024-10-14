package v2

import (
	"context"

	"github.com/coopnorge/go-logger"
)

type Histogram struct {
	metric
}

// Record sends a metric to Statsd.
func (h *Histogram) Record(ctx context.Context, incr float64) {
	if ctx.Err() != nil {
		logger.Errorf("Context was cancelled, will not send metric: %v", ctx.Err())
	}

	if h.metric.client == nil {
		return
	}

	err := h.client.Histogram(h.name, incr, formatTags(h.tags), 1)
	if err != nil {
		logger.WithError(err).Errorf("Failed to Record Histogram: %s", h.name)
	}
}
