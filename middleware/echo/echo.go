package echo

import (
    "fmt"

    "github.com/coopnorge/go-datadog-lib/internal"
    "github.com/coopnorge/go-datadog-lib/tracing"

    "github.com/labstack/echo/v4"
    "gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

// TraceServerMiddleware for Datadog Log Integration, middleware will create span that can be used from context
func TraceServerMiddleware() echo.MiddlewareFunc {
    return func(next echo.HandlerFunc) echo.HandlerFunc {
        return func(c echo.Context) error {
            if c.Request() == nil {
                return fmt.Errorf("unable to extract request from Request Context from Echo it's nil")
            }

            span, spanCtx := tracer.StartSpanFromContext(c.Request().Context(), c.Request().RequestURI, tracer.ResourceName("http.request"))
            defer span.Finish()

            extCtx := internal.ExtendedContextWithMetadata(spanCtx, internal.TraceContextKey{}, tracing.TraceDetails{DatadogSpan: span})
            c.Request().WithContext(extCtx)

            return nil
        }
    }
}
