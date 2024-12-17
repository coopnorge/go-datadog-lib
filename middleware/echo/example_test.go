package echo_test

import (
	"context"

	coopdatadog "github.com/coopnorge/go-datadog-lib/v2"
	coopEchoDatadog "github.com/coopnorge/go-datadog-lib/v2/middleware/echo"
	"github.com/labstack/echo/v4"
)

// `go-datadog-lib` provides middleware for the Echo framework for tracing
// inbound request.
func ExampleTraceServerMiddleware() {
	err := run()
	if err != nil {
		panic(err)
	}
}

func run() error {
	stop, err := coopdatadog.Start(context.Background())
	if err != nil {
		panic(err)
	}
	defer func() {
		err := stop()
		if err != nil {
			panic(err)
		}
	}()

	// ...
	echoServer := echo.New()
	// Some other configuration
	// ...
	// Add middleware to extend context for better traceability
	echoServer.Use(coopEchoDatadog.TraceServerMiddleware())

	return nil
}
