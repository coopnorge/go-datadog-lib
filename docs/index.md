# Go Datadog Library

Reduces the complexity of initializing and using Datadog functionality. See
[Datadog - Getting Started](https://docs.datadoghq.com/getting_started/?site=eu)
for more information about how Datadog works.

## Module documentation

<https://pkg.go.dev/github.com/coopnorge/go-datadog-lib/v2>

## Setup

In order for `go-datadog-lib` to send data to Datadog your service/application
must be configured correctly.

Setting the environment variable `DD_DISABLE` to `true` or any other value that
[`strconv#ParseBool`](https://pkg.go.dev/strconv#ParseBool) can parse to `true`
without returning an error. If `DD_DISABLE` is undefined or a value that
[`strconv#ParseBool`](https://pkg.go.dev/strconv#ParseBool) can parse to `false`
or returns an error the library will be enabled. This is done to ensure that the
library is not disabled in production by accident.

### Kubernetes setup

To instrument an application running inside Kubernetes configure Datadog
[Unified Service Tagging](https://docs.datadoghq.com/getting_started/tagging/unified_service_tagging/?tab=kubernetes)
and set the required environmental variables. If you are using an official Coop
Norge SA Helm chart skip to [application setup](#application-setup).

Kubernetes resource labels:

- `tags.datadoghq.com/service`
- `tags.datadoghq.com/env`
- `tags.datadoghq.com/version`

For resources that defines a template, define the labels the resource template
as well.

Environmental variables:

- `DD_DOGSTATSD_URL`
- `DD_TRACE_AGENT_URL`
- `DD_SERVICE`
- `DD_VERSION`
- `DD_ENV`

> [!NOTE]
>
> Don't forget to set `DD_ENV` for each environment, `production` or `staging`,
> otherwise the application will be not visible in APM Service Catalog.

```yaml title="deployment.yaml"
apiVersion: apps/v1
kind: Deployment
metadata:
  name: my-app
  labels:
    app: my-app
    version: gitc-b24846c73fe50704969ea4bc1e81e3a3a7592296
    tags.datadoghq.com/service: my-app
    tags.datadoghq.com/env: production
    tags.datadoghq.com/version: gitc-b24846c73fe50704969ea4bc1e81e3a3a7592296
spec:
  selector:
    matchLabels:
      app: my-app
  template:
    metadata:
      labels:
        app: my-app
        version: gitc-b24846c73fe50704969ea4bc1e81e3a3a7592296
        tags.datadoghq.com/service: my-app
        tags.datadoghq.com/version: gitc-b24846c73fe50704969ea4bc1e81e3a3a7592296
      annotations:
        proxy.istio.io/config: '{ "holdApplicationUntilProxyStarts": true }'
    spec:
      serviceAccountName: my-app
      containers:
        env:
          - name: DD_DOGSTATSD_URL
            value: "unix:///var/run/datadog/dsd.socket"
          - name: DD_TRACE_AGENT_URL
            value: "unix:///var/run/datadog/apm.socket"
          - name: DD_SERVICE
            valueFrom:
              fieldRef:
                fieldPath: metadata.labels['tags.datadoghq.com/service']
          - name: DD_VERSION
            valueFrom:
              fieldRef:
                fieldPath: metadata.labels['tags.datadoghq.com/version']
          - name: DD_ENV
            valueFrom:
              fieldRef:
                fieldPath: metadata.labels['tags.datadoghq.com/env']
        volumeMounts:
          - name: ddsocket
            mountPath: /var/run/datadog
            readOnly: true
        imagePullPolicy: Always
      volumes:
        - hostPath:
            path: /var/run/datadog/
          name: ddsocket
```

### Application setup

First, we’ll initialize the `go-datadog-lib`. This is required for any
application that exports telemetry.

`coopdatadog.Start` returns a `StopFunc` and an `error`. The `StopFunc` must be
called before the application exits.

```go title="cmd/helloworld/main.go"
package main

import (
	"context"

	coopdatadog "github.com/coopnorge/go-datadog-lib/v2"
)

func main() {
	err := run()
	if err != nil {
		panic(err)
	}
}

func run() error {
	stop, err := coopdatadog.Start(context.Background())
	if err != nil {
		return err
	}
	defer func() {
		err := stop()
		if err != nil {
			panic(err)
		}
	}()

	// ...

	return nil
}
```

> [!NOTE]
>
> After that Datadog will try to connect to the socket and will start to send
> all information in the background.
>
> In different setup, you could have error logs that Datadog cannot connect to
> the socket and tried to connect via HTTP. That could be related to issue when
> your container starts faster and sockets were not ready to communicate with
> Agent or Agent was started later.

## Tracing

### Inbound request tracing

Inbound request can be traced using the gRPC server interceptors or the Echo
HTTP Server middleware. If the upstream application is instrumented to support
distributed tracing the traces will be linked.

#### gRPC server interceptor

`go-datadog-lib` provides gRPC interceptors for tracing inbound request for both
Unary and Stream gRPC endpoints.

```go title="cmd/helloworld/main.go"
package main

import (
	"context"
	"fmt"
	"net"

	coopdatadog "github.com/coopnorge/go-datadog-lib/v2"
	datadogMiddleware "github.com/coopnorge/go-datadog-lib/v2/middleware/grpc"
	"google.golang.org/grpc"
)

func main() {
	err := run()
	if err != nil {
		panic(err)
	}
}

func run() error {
	ctx := context.Background()
	stop, err := coopdatadog.Start(ctx)
	if err != nil {
		panic(err)
	}
	defer func() {
		err := stop()
		if err != nil {
			panic(err)
		}
	}()

	ddOpts := []datadogMiddleware.Option{
		// ...
	}
	serverOpts := []grpc.ServerOption{
		grpc.UnaryInterceptor(datadogMiddleware.UnaryServerInterceptor(ddOpts...)),
		grpc.StreamInterceptor(datadogMiddleware.StreamServerInterceptor(ddOpts...)),
	}

	grpcServer := grpc.NewServer(serverOpts...)

	listener, err := (&net.ListenConfig{}).Listen(ctx, "tcp", "127.0.0.1")
	if err != nil {
		return fmt.Errorf("failed to start tcp listener: %w", err)
	}

	err = grpcServer.Serve(listener)
	if err != nil {
		return nil
	}

	return nil
}
}
```

#### Echo HTTP server middleware

`go-datadog-lib` provides middleware for the Echo framework for tracing inbound
request.

```go title="cmd/helloworld/main.go"
package main

import (
	"context"

	coopdatadog "github.com/coopnorge/go-datadog-lib/v2"
	coopEchoDatadog "github.com/coopnorge/go-datadog-lib/v2/middleware/echo"
	"github.com/labstack/echo/v4"
)

func main() {
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
```

### Outbound request tracing

Outbound requests to other applications/services or storage can be traced. If
the downstream application is instrumented to support distributed tracing the
traces will be linked.

#### gRPC client interceptor

If your application is making gRPC calls, you can add the gRPC
client-interceptor to automatically create child-spans for each gRPC-call. These
spans will also be embedded in the outbound gRPC-metadata, so if you are calling
another service that is also instrumented with Datadog-integration, then you
will enable distributed tracing.

It is important that the context used in the RPC contains trace-information,
preferably created from any server middleware from this module.

```go
package main

import (
	"context"

	"github.com/DataDog/dd-trace-go/v2/ddtrace/tracer"
	datadogMiddleware "github.com/coopnorge/go-datadog-lib/v2/middleware/grpc"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/testing/testpb"
	"google.golang.org/grpc"
)

func main() {
	unary()
	stream()
}

func unary() {
	ctx := context.Background()

	cc, err := grpc.NewClient(
		"https://example.com",
		grpc.WithUnaryInterceptor(datadogMiddleware.UnaryClientInterceptor()),
	)
	if err != nil {
		panic(err)
	}
	defer cc.Close()

	client := testpb.NewTestServiceClient(cc)

	span, ctx := tracer.StartSpanFromContext(ctx, "grpc.request")
	resp, err := client.Ping(ctx, &testpb.PingRequest{})
	span.Finish(tracer.WithError(err))
	if err != nil {
		panic(err)
	}
	println(resp)
}

func stream() {
	ctx := context.Background()

	cc, err := grpc.NewClient(
		"https://example.com",
		grpc.WithStreamInterceptor(datadogMiddleware.StreamClientInterceptor()),
	)
	if err != nil {
		panic(err)
	}
	defer cc.Close()

	client := testpb.NewTestServiceClient(cc)

	span, ctx := tracer.StartSpanFromContext(ctx, "grpc.stream")
	stream, err := client.PingStream(ctx)
	defer span.Finish()
	if err != nil {
		span.Finish(tracer.WithError(err))
		panic(err)
	}
	resp, err := stream.Recv()
	if err != nil {
		span.Finish(tracer.WithError(err))
		panic(err)
	}
	println(resp)
}
```

#### HTTP client middleware

If your application is making HTTP calls, you can add the HTTP
client-interceptor to automatically create child-spans for each HTTP-call. These
spans will also be embedded in the outbound HTTP Headers, so if you are calling
another service that is also instrumented with Datadog-integration, then you
will enable distributed tracing.

It is important that the context used in the `http.Request` contains
trace-information, preferably created from any server middleware from this
module.

```go
package main

import (
	"context"
	"net/http"
	"time"

	"github.com/DataDog/dd-trace-go/v2/ddtrace/tracer"
	datadogMiddleware "github.com/coopnorge/go-datadog-lib/v2/middleware/http"
)

func main() {
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
```

#### Standard library SQL middleware

If your application is making calls to a database, you can add the database
driver to automatically create child-spans for each database-call.

It is important that the context used in the call to the database contains
trace-information, preferably created from any server middleware from this
module.

```go
package main

import (
	"context"

	"github.com/DataDog/dd-trace-go/v2/ddtrace/tracer"
	ddDatabase "github.com/coopnorge/go-datadog-lib/v2/middleware/database"
	mysqlDriver "github.com/go-sql-driver/mysql"
)

func main() {
	ctx := context.Background()

	dsn := "example.com/users"
	db, err := ddDatabase.RegisterDriverAndOpen("mysql", mysqlDriver.MySQLDriver{}, dsn)
	if err != nil {
		panic(err)
	}

	span, ctx := tracer.StartSpanFromContext(ctx, "http.request")
	defer span.Finish()
	rows, err := db.QueryContext(ctx, "SELECT * FROM users")
	if err != nil {
		span.Finish(tracer.WithError(err))
		panic(err)
	}
	println(rows)
}
```

#### GORM middleware

If your application is using GORM to make calls to a database, you can add the
GORM middleware to automatically create child-spans for each database-call.

It is important that the context used in the call to the database contains
trace-information, preferably created from any server middleware from this
module.

```go
package main

import (
	"context"

	ddDatabase "github.com/coopnorge/go-datadog-lib/v2/middleware/database"
	ddGorm "github.com/coopnorge/go-datadog-lib/v2/middleware/gorm"
	mysqlDriver "github.com/go-sql-driver/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type User struct{}

func main() {
	ctx := context.Background()

	dsn := "example.com/users"
	gormDB, err := ddGorm.NewORM(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	user := &User{}
	tx := gormDB.WithContext(ctx).Select("*").First(user)

	println(tx)
}
```

You can also combine this with the standard library tracer-middleware:

```go
package main

import (
	"context"

	ddDatabase "github.com/coopnorge/go-datadog-lib/v2/middleware/database"
	ddGorm "github.com/coopnorge/go-datadog-lib/v2/middleware/gorm"
	mysqlDriver "github.com/go-sql-driver/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type User struct{}

func main() {
	ctx := context.Background()

	dsn := "example.com/users"
	db, err := ddDatabase.RegisterDriverAndOpen("mysql", mysqlDriver.MySQLDriver{}, dsn)
	if err != nil {
		panic(err)
	}

	gormDB, err := ddGorm.NewORM(mysql.New(mysql.Config{Conn: db}), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	user := &User{}
	tx := gormDB.WithContext(ctx).Select("*").First(user)

	println(tx)
}
```

## Metrics

The package `github.com/coopnorge/go-datadog-lib/v2/metrics` contains function
to send telemetry. Before sending metrics `go-datadog-lib` has to be
[initialized](#application-setup).

Use cases:

- track value of carts in an e-commerce shop.
- count failed authentication attempts.

```go
package main

import (
	"context"

	coopdatadog "github.com/coopnorge/go-datadog-lib/v2"
	"github.com/coopnorge/go-datadog-lib/v2/metrics"
)

func main() {
	err := run()
	if err != nil {
		panic(err)
	}
}

func run() error {
	stop, err := coopdatadog.Start(context.Background())
	if err != nil {
		return err
	}
	defer func() {
		err := stop()
		if err != nil {
			panic(err)
		}
	}()

	metrics.Incr("my-metric")

	return nil
}
```

### Example how to send metrics

When you have `BaseMetricCollector` from pkg `metrics` you can call create
records in Datadog.

```go
package my_metric

import (
	"context"

	"github.com/coopnorge/go-datadog-lib/v2/config"
	"github.com/coopnorge/go-datadog-lib/v2/metric"
)

func Example()  {
	ddClient, ddClientErr := metrics.NewDatadogMetrics(new(config.DatadogConfig))
	if ddClient != nil {
    // Handle error / log error
  }
	ddMetricCollector := metrics.NewBaseMetricCollector(ddClient)

	tMetricData := metrics.Data{
		Name:  "RuntimeTest",
		Type:  metrics.MetricTypeEvent,
		Value: float64(42),
		MetricTags: []metrics.Tag{
			{Name: "Show", Value: "Case"},
		},
	}

	ddMetricCollector.AddMetric(context.Background(), tMetricData)
}
```

## Datadog Context Log Hook

Relate log-entries to traces in Datadog. Configure
`github.com/cooopnorge/go-logger` with a Hook (documentation:
[Inventory](https://inventory.internal.coop/docs/default/component/go-logger/#hooks),
[GitHub](https://github.com/coopnorge/go-logger/blob/main/docs/index.md#hooks))
to capture the `trace_id` `span_id`.

```go
package main

import (
	"context"

	"github.com/coopnorge/go-datadog-lib/v2/tracelogger"
	"github.com/coopnorge/go-logger"

	"github.com/DataDog/dd-trace-go/v2/ddtrace/tracer"
)

func main() {
	logger.ConfigureGlobalLogger(
		logger.WithHook(tracelogger.NewHook()),
	)
	ctx := context.Background()
	a(ctx)
}

func a(ctx context.Context) {
	span, ctx := tracer.StartSpanFromContext(ctx, "a")
	err := b(ctx)
	span.Finish(tracer.WithError(err))
}

func b(ctx context.Context) error {
	logger.WithContext(ctx).Info("Hello")
	return nil
}

// Output:
// {"dd.span_id":8047616890857967865,"dd.trace_id":8160264448608745330,"file":"/srv/workspace/app/main.go:25","function":"github.com/coopnorge/app/main.b","level":"info","msg":"Hello","time":"2024-09-12T19:01:34+02:00"}
```
