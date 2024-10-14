package v2

import (
	"context"

	"github.com/coopnorge/go-logger"
)

type Counter struct {
	metric
}

func (c *Counter) Count(ctx context.Context, incr int64, tags ...Tag) {
	if ctx.Err() != nil {
		logger.Errorf("Context was cancelled, will not send metric: %v", ctx.Err())
	}

	if c.metric.client == nil {
		return
	}
	err := c.client.Count(c.name, incr, formatTags(append(c.tags, tags...)), 1)
	if err != nil {
		logger.WithError(err).Errorf("Failed to Send Count on Counter: %s", c.name)
	}
}
