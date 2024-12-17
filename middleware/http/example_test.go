package http_test

import (
	"context"
	"net/http"
	"time"

	datadogMiddleware "github.com/coopnorge/go-datadog-lib/v2/middleware/http"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

func ExampleAddTracingToClient() {
	ctx := context.Background()

	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	client = datadogMiddleware.AddTracingToClient(client)

	span, ctx := tracer.StartSpanFromContext(ctx, "http.request")
	defer span.Finish()
	req, err := http.NewRequestWithContext(ctx, "GET", "https://example.com", nil)
	if err != nil {
		span.Finish(tracer.WithError(err))
		panic(err)
	}
	resp, err := client.Do(req)
	if err != nil {
		span.Finish(tracer.WithError(err))
		panic(err)
	}
	println(resp)
}
