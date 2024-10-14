package v2

import (
	"fmt"
	"strings"

	"github.com/DataDog/datadog-go/v5/statsd"
	"github.com/iancoleman/strcase"
)

type Tag struct {
	Key   string
	Value string
}

type metric struct {
	name        string
	description string
	unit        string
	client      *statsd.Client
	tags        []Tag
}

func formatTags(tags []Tag) (res []string) {
	for _, tag := range tags {
		tagName := strings.ToLower(strcase.ToKebab(tag.Key))
		res = append(res, fmt.Sprintf("%s:%s", tagName, tag.Value))
	}
	return
}

type MetricOption func(g *metric)

func WithTags(tags ...Tag) MetricOption {
	return func(g *metric) {
		g.tags = append(g.tags, tags...)
	}
}

func WithDescription(description string) MetricOption {
	return func(g *metric) {
		g.description = description
	}
}

func WithUnit(unit string) MetricOption {
	return func(g *metric) {
		g.unit = unit
	}
}

func WithTag(key, value string) MetricOption {
	return WithTags(Tag{Key: key, Value: value})
}

type StatsdProvider struct {
	client *statsd.Client
}

type Provider interface {
	NewCounter(name string, options ...MetricOption) *Counter
	NewGauge(name string, options ...MetricOption) *Gauge
	NewHistogram(name string, options ...MetricOption) *Histogram
}

func NewStatsdProvider(client *statsd.Client) *StatsdProvider {
	return &StatsdProvider{client: client}
}

func (p *StatsdProvider) NewCounter(name string, options ...MetricOption) *Counter {
	c := Counter{
		metric: metric{
			name:   name,
			client: p.client,
		},
	}
	for _, opt := range options {
		opt(&c.metric)
	}
	return &c
}

func (p *StatsdProvider) NewGauge(name string, options ...MetricOption) *Gauge {
	g := Gauge{
		metric: metric{
			name:   name,
			client: p.client,
		},
	}
	for _, opt := range options {
		opt(&g.metric)
	}
	return &g
}

func (p *StatsdProvider) NewHistogram(name string, options ...MetricOption) *Histogram {
	h := Histogram{
		metric: metric{
			name:   name,
			client: p.client,
		},
	}
	for _, opt := range options {
		opt(&h.metric)
	}
	return &h
}

type NoopProvider struct{}

func (p *NoopProvider) NewCounter(name string, options ...MetricOption) *Counter {
	metric := Counter{
		metric: metric{
			name: name,
		},
	}
	for _, opt := range options {
		opt(&metric.metric)
	}
	return &metric
}

func (p *NoopProvider) NewGauge(name string, options ...MetricOption) *Gauge {
	metric := Gauge{
		metric: metric{
			name: name,
		},
	}
	for _, opt := range options {
		opt(&metric.metric)
	}
	return &metric
}

func (p *NoopProvider) NewHistogram(name string, options ...MetricOption) *Histogram {
	metric := Histogram{
		metric: metric{
			name: name,
		},
	}
	for _, opt := range options {
		opt(&metric.metric)
	}
	return &metric
}
