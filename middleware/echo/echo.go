package echo

import (
	"github.com/coopnorge/go-datadog-lib/v2/internal"
	"github.com/labstack/echo/v4"
	ddEcho "gopkg.in/DataDog/dd-trace-go.v1/contrib/labstack/echo.v4"
)

// TraceServerMiddleware for Datadog Log Integration, middleware will create span that can be used from context
func TraceServerMiddleware() echo.MiddlewareFunc {
	if internal.IsDatadogDisabled() {
		return noOpMiddlewareFunc()
	}

	return ddEcho.Middleware()
}

func noOpMiddlewareFunc() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			return next(c)
		}
	}
}
