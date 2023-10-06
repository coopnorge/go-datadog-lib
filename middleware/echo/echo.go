package echo

import (
	"fmt"

	"github.com/coopnorge/go-datadog-lib/v2/internal"
	"github.com/coopnorge/go-datadog-lib/v2/tracing"

	"github.com/labstack/echo/v4"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

// TraceServerMiddleware for Datadog Log Integration, middleware will create span that can be used from context
func TraceServerMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			req := c.Request()
			if req == nil {
				return fmt.Errorf("unable to extract request from Echo Request Context, returned nil")
			}

			opts := []ddtrace.StartSpanOption{tracer.Measured()}
			if spanCtx, err := tracer.Extract(tracer.HTTPHeadersCarrier(req.Header)); err == nil {
				opts = append(opts, tracer.ChildOf(spanCtx))
			}

			span, spanCtx := tracer.StartSpanFromContext(req.Context(), req.RequestURI, opts...)
			defer span.Finish()

			extCtx := internal.ExtendedContextWithMetadata(
				spanCtx,
				internal.TraceContextKey{},
				tracing.TraceDetails{DatadogSpan: span},
			)

			c.SetRequest(req.WithContext(extCtx))

			return next(c)
		}
	}
}

// TraceServerMiddlewareExperimental is experimental, and will be removed in the next non-pre-release version. Used for testing a new way of setting up the middleware.
func TraceServerMiddlewareExperimental() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			req := c.Request()
			if req == nil {
				return fmt.Errorf("unable to extract request from Echo Request Context, returned nil")
			}

			opts := []ddtrace.StartSpanOption{tracer.Measured()}
			if spanCtx, err := tracer.Extract(tracer.HTTPHeadersCarrier(req.Header)); err == nil {
				opts = append(opts, tracer.ChildOf(spanCtx))
			}

			span, spanCtx := tracer.StartSpanFromContext(req.Context(), req.RequestURI, opts...)
			defer span.Finish()

			extCtx := internal.ExtendedContextWithMetadata(
				spanCtx,
				internal.TraceContextKey{},
				tracing.TraceDetails{DatadogSpan: span},
			)

			c.SetRequest(req.WithContext(extCtx))

			return next(c)
		}
	}
}
