package database

import (
	"database/sql"
	"database/sql/driver"
	"os"

	sqltrace "github.com/DataDog/dd-trace-go/contrib/database/sql/v2"
	"github.com/coopnorge/go-datadog-lib/v2/internal"
)

// RegisterDriverAndOpen registers the selected driver with the datadog-lib, and opens a connection to the database using the dsn.
func RegisterDriverAndOpen(driverName string, driver driver.Driver, dsn string, options ...Option) (*sql.DB, error) {
	if internal.IsDatadogDisabled() {
		sql.Register(driverName, driver)
		return sql.Open(driverName, dsn)
	}

	cfg := defaults()
	for _, opt := range options {
		opt(cfg)
	}
	opts := make([]sqltrace.Option, 0, 3+len(cfg.tags))
	if cfg.serviceName != "" {
		opts = append(opts, sqltrace.WithService(cfg.serviceName))
	}
	if cfg.childSpansOnly {
		opts = append(opts, sqltrace.WithChildSpansOnly())
	}
	for k, v := range cfg.tags {
		opts = append(opts, sqltrace.WithCustomTag(k, v))
	}
	if len(cfg.ignoredQueryTypes) > 0 {
		typed := make([]sqltrace.QueryType, 0, len(cfg.ignoredQueryTypes))
		for i := range cfg.ignoredQueryTypes {
			typed = append(typed, sqltrace.QueryType(cfg.ignoredQueryTypes[i]))
		}
		opts = append(opts, sqltrace.WithIgnoreQueryTypes(typed...))
	}

	sqltrace.Register(driverName, driver, opts...)
	return sqltrace.Open(driverName, dsn)
}

type config struct {
	serviceName       string
	childSpansOnly    bool
	tags              map[string]interface{}
	ignoredQueryTypes []string
}

func defaults() *config {
	serviceName := os.Getenv("DD_SERVICE")
	return &config{
		serviceName:    serviceName,
		childSpansOnly: true,
		tags:           nil,
		ignoredQueryTypes: []string{
			string(sqltrace.QueryTypeConnect),
			string(sqltrace.QueryTypePing),
			string(sqltrace.QueryTypePrepare),
			string(sqltrace.QueryTypeClose),
		},
	}
}

// Option allows for overriding our default-config.
type Option func(cfg *config)

// WithServiceName overrides the service-name set in environment-variable "DD_SERVICE".
func WithServiceName(serviceName string) Option {
	return func(cfg *config) {
		cfg.serviceName = serviceName
	}
}

// WithCustomTag will attach the value to the span tagged by the key.
func WithCustomTag(key string, value interface{}) Option {
	return func(cfg *config) {
		if cfg.tags == nil {
			cfg.tags = make(map[string]interface{})
		}
		cfg.tags[key] = value
	}
}

// WithChildSpansOnly causes spans to be created only when there is an existing parent span in the Context.
func WithChildSpansOnly(childSpansOnly bool) Option {
	return func(cfg *config) {
		cfg.childSpansOnly = childSpansOnly
	}
}

// WithIgnoreQueryTypes specifies the query types for which spans should not be created.
// Will replace any existing ignored query-types, so it must be an exhaustive list.
// See available QueryTypes here: https://pkg.go.dev/github.com/DataDog/dd-trace-go/contrib/database/sql/v2#pkg-constants
func WithIgnoreQueryTypes(ignoredQueryTypes ...string) Option {
	return func(cfg *config) {
		cfg.ignoredQueryTypes = ignoredQueryTypes
	}
}
