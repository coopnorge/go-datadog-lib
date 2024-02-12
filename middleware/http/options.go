package http

import (
	"net/http"
	"os"

	httptrace "gopkg.in/DataDog/dd-trace-go.v1/contrib/net/http"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

// These options, unless otherwise specified, will be valid for both client and server interceptors.

// ResourceNamer is a function that  will be called to determine what the ResourceName of the span should be.
type ResourceNamer func(req *http.Request) string

// RequestIgnorer is a function that  will be called to determine if the request should be traced or not.
type RequestIgnorer func(req *http.Request) bool

type config struct {
	serviceName    string
	resourceNamer  ResourceNamer
	requestIgnorer RequestIgnorer
	tags           map[string]any
}

func defaults() *config {
	serviceName := os.Getenv("DD_SERVICE")
	return &config{
		serviceName:    serviceName,
		resourceNamer:  FullURLResourceNamer(),
		requestIgnorer: nil,
		tags:           nil,
	}
}

// convertClientOptions converts the Options to httptrace-typed options.
// httptrace-clients use httptrace.RoundTripperOption.
// httptrace-servers use httptrace.Option.
func convertClientOptions(options ...Option) []httptrace.RoundTripperOption {
	cfg := defaults()
	for _, opt := range options {
		opt(cfg)
	}
	opts := make([]httptrace.RoundTripperOption, 0, 3+len(cfg.tags))
	if cfg.serviceName != "" {
		opts = append(opts, httptrace.RTWithServiceName(cfg.serviceName))
	}
	if cfg.resourceNamer != nil {
		opts = append(opts, httptrace.RTWithResourceNamer(cfg.resourceNamer))
	}
	if cfg.requestIgnorer != nil {
		opts = append(opts, httptrace.RTWithIgnoreRequest(cfg.requestIgnorer))
	}
	for k, v := range cfg.tags {
		opts = append(opts, httptrace.RTWithSpanOptions(tracer.Tag(k, v)))
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

// WithResourceNamer specifies a function that will be called to determine what
// the ResourceName of the span should be.
// For example, in the span "http.client https://coop.no/", the ResourceName is "https://coop.no/"
func WithResourceNamer(resourceNamer ResourceNamer) Option {
	return func(cfg *config) {
		cfg.resourceNamer = resourceNamer
	}
}

// WithRequestIgnorer specifies a function that will be called to determined if
// the request should be traced or not.
func WithRequestIgnorer(requestIgnorer RequestIgnorer) Option {
	return func(cfg *config) {
		cfg.requestIgnorer = requestIgnorer
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

// StaticResourceNamer will set every span's ResourceName to str.
func StaticResourceNamer(str string) ResourceNamer {
	return func(_ *http.Request) string {
		return str
	}
}

// FullURLWithParamsResourceNamer will name the ResourceName to "GET https://www.coop.no/api/some-service/some-endpoint?foo=bar"
// NOTE! This might leak unintended user-ids, credentials, or other things that are part of a URL.
func FullURLWithParamsResourceNamer() ResourceNamer {
	return func(req *http.Request) string {
		return req.Method + " " + req.URL.Redacted()
	}
}

// FullURLResourceNamer will name the ResourceName to "GET https://www.coop.no/api/some-service/some-endpoint"
// NOTE! This might leak unintended user-ids, if they are part of the URL-path.
func FullURLResourceNamer() ResourceNamer {
	return func(req *http.Request) string {
		u := req.URL
		return req.Method + " " + u.Scheme + "://" + u.Host + u.Path
	}
}

// HostResourceNamer will name the ResourceName to "www.coop.no"
func HostResourceNamer() ResourceNamer {
	return func(req *http.Request) string {
		return req.URL.Host
	}
}
