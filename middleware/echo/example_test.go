package echo_test

import (
	"context"

	coopdatadog "github.com/coopnorge/go-datadog-lib/v2"
	coopEchoDatadog "github.com/coopnorge/go-datadog-lib/v2/middleware/echo"
	"github.com/labstack/echo/v4"
)

func ExampleWrap() {
	err := runWrap()
	if err != nil {
		panic(err)
	}
}

func runWrap() error {
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
	// Wrap the Echo server to extend context for better traceability
	coopEchoDatadog.Wrap(echoServer)

	return nil
}

// `go-datadog-lib` provides middleware for the Echo framework for tracing
// inbound request.
func ExampleTraceServerMiddleware() {
	err := runMiddleware()
	if err != nil {
		panic(err)
	}
}

func runMiddleware() error {
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
