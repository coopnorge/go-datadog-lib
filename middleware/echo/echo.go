package echo

import (
	ddEcho "github.com/DataDog/dd-trace-go/contrib/labstack/echo.v4/v2"
	"github.com/coopnorge/go-datadog-lib/v2/internal"
	"github.com/labstack/echo/v5"
)

// Wrap configures the provided [echo.Echo] and returns it. This is a
// convenience function that calls [ddEcho.Wrap] with the provided [echo.Echo]
// It is recommended to use this if you want to benefit from future tracer
// features that require additional properties to be configured without having
// to update your code.
func Wrap(e *echo.Echo) *echo.Echo {
	if internal.IsDatadogDisabled() {
		return e
	}

	e = ddEcho.Wrap(e)
	return e
}

// TraceServerMiddleware for Datadog Log Integration, middleware will create span that can be used from context
//
// Deprecated: Use of the [Wrap] function is recommended instead of directly calling
// [echo.TraceServerMiddleware], as [Wrap] activates all available features automatically.
func TraceServerMiddleware() echo.MiddlewareFunc {
	if internal.IsDatadogDisabled() {
		return noOpMiddlewareFunc()
	}

	return ddEcho.Middleware() // nolint:staticcheck
}

func noOpMiddlewareFunc() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			return next(c)
		}
	}
}
