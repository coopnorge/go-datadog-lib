package echo

import (
	"github.com/labstack/echo/v4"
	ddEcho "gopkg.in/DataDog/dd-trace-go.v1/contrib/labstack/echo.v4"
)

// TraceServerMiddleware for Datadog Log Integration, middleware will create span that can be used from context
func TraceServerMiddleware() echo.MiddlewareFunc {
	return ddEcho.Middleware()
}
