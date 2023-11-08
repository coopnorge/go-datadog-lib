package gorm

import (
	"database/sql"
	"os"

	gormtrace "gopkg.in/DataDog/dd-trace-go.v1/contrib/gorm.io/gorm.v1"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// NewMySQLORM returns a new gorm DB instance.
func NewMySQLORM(db *sql.DB, gormCfg *gorm.Config, options ...Option) (*gorm.DB, error) {
	cfg := defaults()
	for _, opt := range options {
		opt(cfg)
	}
	opts := make([]gormtrace.Option, 0, 2)
	if cfg.serviceName != "" {
		opts = append(opts, gormtrace.WithServiceName(cfg.serviceName))
	}
	for k, v := range cfg.tags {
		v := v
		staticTagger := func(db *gorm.DB) any {
			return v
		}
		opts = append(opts, gormtrace.WithCustomTag(k, staticTagger))
	}

	dialector := mysql.New(mysql.Config{Conn: db})
	return gormtrace.Open(dialector, gormCfg, opts...)
}

type config struct {
	serviceName string
	tags        map[string]interface{}
}

func defaults() *config {
	serviceName := os.Getenv("DD_SERVICE")
	return &config{
		serviceName: serviceName,
		tags:        nil,
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
