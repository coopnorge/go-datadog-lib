package grpc

import (
	"os"

	"google.golang.org/grpc/codes"
	ddGrpc "gopkg.in/DataDog/dd-trace-go.v1/contrib/google.golang.org/grpc"
)

// These options, unless otherwise specified, will be valid for both client and server interceptors.

type config struct {
	serviceName     string
	nonErrorCodes   []codes.Code
	untracedMethods []string
	tags            map[string]any
}

func defaults() *config {
	serviceName := os.Getenv("DD_SERVICE")
	return &config{
		serviceName:     serviceName,
		nonErrorCodes:   []codes.Code{codes.Canceled},
		untracedMethods: nil,
		tags:            nil,
	}
}

func convertOptions(options ...Option) []ddGrpc.Option {
	cfg := defaults()
	for _, opt := range options {
		opt(cfg)
	}
	opts := make([]ddGrpc.Option, 0, 3+len(cfg.tags))
	if cfg.serviceName != "" {
		opts = append(opts, ddGrpc.WithServiceName(cfg.serviceName))
	}
	if len(cfg.nonErrorCodes) > 0 {
		opts = append(opts, ddGrpc.NonErrorCodes(cfg.nonErrorCodes...))
	}
	if len(cfg.untracedMethods) > 0 {
		opts = append(opts, ddGrpc.WithUntracedMethods(cfg.untracedMethods...))
	}
	for k, v := range cfg.tags {
		opts = append(opts, ddGrpc.WithCustomTag(k, v))
	}
	return opts
}

// Option allows for overriding our default-config.
type Option func(cfg *config)

// WithServiceName overrides the service-name set in environment-variable "DD_SERVICE".
func WithServiceName(serviceName string) Option {
	return func(cfg *config) {
		cfg.serviceName = serviceName
	}
}

// WithNonErrorCodes determines the list of codes which will not be considered errors in instrumentation.
// This call overrides the default handling of codes.Canceled as a non-error.
func WithNonErrorCodes(cs ...codes.Code) Option {
	return func(cfg *config) {
		cfg.nonErrorCodes = cs
	}
}

// WithUntracedMethods specifies full methods to be ignored by the server side and client
// side interceptors. When a request's full method is in 'methods', no spans will be created.
func WithUntracedMethods(methods ...string) Option {
	return func(cfg *config) {
		cfg.untracedMethods = methods
	}
}

// WithCustomTag will attach the value to the span tagged by the key.
func WithCustomTag(key string, value any) Option {
	return func(cfg *config) {
		if cfg.tags == nil {
			cfg.tags = make(map[string]any)
		}
		cfg.tags[key] = value
	}
}
