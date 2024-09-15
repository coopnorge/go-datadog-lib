# Go Datadog Library

Reduces the complexity of initializing these services inside your application.
Also provides abstract code to work with metrics.

## APM

In Coop Norge, our default setup is tracing applications with CPU profiling
support, that is enabled by default in the package bootstrap.

## Custom metrics - DD StatsD

Inside the `metric` package you can find the base client for Datadog StatsD and
simple metric service that allows sending `Incr`, `Gauge`, and `Count`.

## How Datadog works

![Datadog diagram](dd_com_app.png)

## Setup

In order for `go-datadog-lib` to send data to Datadog your service/application
must be configured correctly.

Setting the environment variable `DD_DISABLE` to `true` or any other value that
[`strconv#ParseBool`](https://pkg.go.dev/strconv#ParseBool) can parse to `true`
without returning an error. If `DD_DISABLE` is undefined or a value that
[`strconv#ParseBool`](https://pkg.go.dev/strconv#ParseBool) can parse to
`false` or returns an error the library will be enabled. This is done to ensure
that the library is not disabled in production by accident.

### 1. Setup container

Following configuration example related to Kubernetes and Kustomize.

NOTE: Don't forget to set `DD_ENV` for each environment, otherwise it will be
not visible in APM list.

```yaml
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
              value: "/var/run/datadog/apm.socket"
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

It's how the application will be shown in Datadog APM.

```yaml
- name: DD_SERVICE
  valueFrom:
    fieldRef:
      fieldPath: metadata.labels['tags.datadoghq.com/service']
```

Depending on your policy,
you can have an API version or tag/commit from git.

```yaml
- name: DD_VERSION
  valueFrom:
    fieldRef:
      fieldPath: metadata.labels['tags.datadoghq.com/version']
```

### 2. Application setup

Create pkg configuration for bootstrap Datadog.

```go
package main

import (
	coopdatadog "github.com/coopnorge/go-datadog-lib/v2"
	"github.com/coopnorge/go-datadog-lib/v2/config"
)

func main() {
	// Your app initialization
	/// ... 
	// From your core configuration add datadog related values
	ddCfg := config.LoadDatadogConfigFromEnvVars()

	// When you start other processes start datadog
	startDatadogServiceError := coopdatadog.StartDatadog(ddCfg, coopdatadog.ConnectionTypeSocket)
	if startDatadogServiceError != nil {
	// Handle error / log error
	}

	// Stop datadog with yours other processes
	handleGracefulShutdown(coopdatadog.GracefulDatadogShutdown)
	// or simply call with defer
	defer coopdatadog.GracefulDatadogShutdown()
}
```

### 3. Middleware gRPC server

To have better tracing you need add to your gRPC custom middleware that will
extend context.

It's needed to relate logs with your trace data in APM.

To do that, simply add Go - Datadog middleware to your gRPC interceptor.

Take a look at the function `UnaryServerInterceptor` in
[`github.com/coopnorge/go-datadog-lib/blob/main/middleware/grpc/server.go`](https://github.com/coopnorge/go-datadog-lib/blob/main/middleware/grpc/server.go).

```go
import (
	coopdatadog "github.com/coopnorge/go-datadog-lib/v2"
	datadogMiddleware "github.com/coopnorge/go-datadog-lib/v2/middleware/grpc"
	"google.golang.org/grpc"
)

func main() {
	err := coopdatadog.StartDatadog(...)
	if err != nil {
		panic(err)
	}
	defer coopdatadog.GracefulDatadogShutdown()


	ddOpts := []datadogMiddleware.Option{
		// ...
	}
	serverOpts := []grpc.ServerOption{
		grpc.UnaryInterceptor(datadogMiddleware.UnaryServerInterceptor(ddOpts...)),
		grpc.StreamInterceptor(datadogMiddleware.StreamServerInterceptor(ddOpts...))
	}

	grpcServer := grpc.NewServer(serverOpts...)
}
```

### 3. Middleware echo server

Same as gRPC middleware but for Echo framework. It will extend request context
and will allow to create nested spans for it, also correlate with logs.

Example:

```go
package myServer

import (
	coopEchoDatadog "github.com/coopnorge/go-datadog-lib/v2/middleware/echo"
	"github.com/labstack/echo/v4"
)

func MyServer() {
	err := coopdatadog.StartDatadog(...)
	if err != nil {
		panic(err)
	}
	defer coopdatadog.GracefulDatadogShutdown()

	// ...
	echoServer := echo.New()
	// Some other configuration
	// ...
	// Add middleware to extend context for better traceability
	echoServer.Use(coopEchoDatadog.TraceServerMiddleware())
}
```

#### Common issue

After that Datadog will try to connect to the socket and will start to send all
information in the background.

In different setup, you could have error logs that Datadog cannot connect to
the socket and tried to connect via HTTP. That could be related to issue when
your container starts faster and sockets were not ready to communicate with
Agent or Agent was started later.

### 4. Middleware gRPC client

If your application is making outgoing gRPC calls, you can add the gRPC
client-interceptor to automatically create child-spans for each outgoing gRPC-call.
These spans will also be embedded in the outgoing gRPC-metadata, so if you are calling
another service that is also instrumented with Datadog-integration, then you will
enable distributed tracing.

It is important that the context used in the RPC contains trace-information,
preferably created from any server middleware from this module.

Example:

```go
import (
	datadogMiddleware "github.com/coopnorge/go-datadog-lib/v2/middleware/grpc"
	"google.golang.org/grpc"
	myprotov1 "some/import/path"
)

func foo() {
	cc, err := grpc.Dial(
		url,
		grpc.WithUnaryInterceptor(ddGrpc.UnaryClientInterceptor()),
	)
	if err != nil {
		panic(err)
	}

	client := myprotov1.NewFoobarAPIClient(cc)

	_, err := client.SomeUnaryRPC(ctx, myprotov1.SomeUnaryRPCRequest{})
	if err != nil {
		panic(err)
	}
}
```

### 4. Middleware HTTP client

If your application is making outgoing HTTP calls, you can add the HTTP
client-interceptor to automatically create child-spans for each outgoing HTTP-call.
These spans will also be embedded in the outgoing HTTP Headers, so if you are calling
another service that is also instrumented with Datadog-integration, then you will
enable distributed tracing.

It is important that the context used in the `http.Request` contains
trace-information, preferably created from any server middleware from this
module.

Example:

```go
import (
	datadogMiddleware "github.com/coopnorge/go-datadog-lib/v2/middleware/grpc"
	"google.golang.org/grpc"
)

func foo() {
	client := &http.Client{
		Timeout: 10*time.Second,
	}
	client = datadogMiddleware.AddTracingToClient(client)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		panic(err)
	}

	_, err := client.Do(req)
	if err != nil {
		panic(err)
	}
}
```

### 4. Middleware database standard library

If your application is making outgoing call to a database, you can add the database
driver to automatically create child-spans for each outgoing database-call.

It is important that the context used in the call to the database contains trace-information,
preferably created from any server middleware from this module.

Example:

```go
import (
	ddDatabase "github.com/coopnorge/go-datadog-lib/v2/middleware/database"
	mysqlDriver "github.com/go-sql-driver/mysql"
)

func foo() {
	// Example using mysql driver
	db, err := ddDatabase.RegisterDriverAndOpen("mysql", mysqlDriver.MySQLDriver{}, dsn, opts...)
	if err != nil{
		panic(err)
	}

	_, err := db.QueryContext(ctx, "SELECT * FROM users")
}
```

### 4. Middleware database GORM

If your application is using GORM to make outgoing calls to a database, you can
add the GORM middleware to automatically create child-spans for each outgoing database-call.

It is important that the context used in the call to the database contains trace-information,
preferably created from any server middleware from this module.

Example:

```go
import (
	ddGorm "github.com/coopnorge/go-datadog-lib/v2/middleware/gorm"
	mysqlDriver "github.com/go-sql-driver/mysql"
)

func foo() {
	// Example using mysql driver
	gormDB, err := ddGorm.NewORM(mysql.New(mysql.Config{Conn: sqlDb}), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	user := &entity.User{}
	_ := gormDB.WithContext(ctx).Select("*").First(user)
}
```

You can also combine this with the standard library tracer-middleware:

```go
import (
	ddDatabase "github.com/coopnorge/go-datadog-lib/v2/middleware/database"
	mysqlDriver "github.com/go-sql-driver/mysql"
)

func foo() {
	// Example using mysql driver
	db, err := ddDatabase.RegisterDriverAndOpen("mysql", mysqlDriver.MySQLDriver{}, dsn, opts...)
	if err != nil{
		panic(err)
	}

	gormDB, err := ddGorm.NewORM(mysql.New(mysql.Config{Conn: db}), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	user := &entity.User{}
	_ := gormDB.WithContext(ctx).Select("*").First(user)
}
```

## Metric - Datadog StatsD

Datadog supports custom metrics that you can utilize depending on the
application.

For example, you could use it to track value of cart in side e-commerce shop.

Or you could register events for auth attempts.

All depends on the case and what you're looking forward to achieve.

### How to use StatsD in Go - Datadog

There is created an abstract client that simply connects to Datadog StatsD
service.

Also, you will already implement simple metric the collector that you can
extend or just use it to send your events and measurements.

### How to initialize Go - Datadog

To prepare the configuration you need to look at the `Setup` section. After
that you can create new instance of Datadog client for DD StatsD.

```go
package your_pkg

import (
	"github.com/coopnorge/go-datadog-lib/v2/config"
	"github.com/coopnorge/go-datadog-lib/v2/metric"
)

func MyServiceContainer(ddCfg *config.DatadogConfig) error {
	// After that you will have pure DD StatsD client
	ddClient := metrics.NewDatadogMetrics(ddCfg)

	// If you need simple metric collector then create
	ddMetricCollector, ddMetricCollectorErr := metrics.NewBaseMetricCollector(ddClient)
	if ddMetricCollectorErr != nil {
		// Handle error / log error
  }
	// ddMetricCollector -> *BaseMetricCollector allows you to send metrics to Datadog
    
	// ensure the metrics are sent before the program is terminated
	defer ddMetricCollector.GracefulShutdown()
}
```

### Example how to send metrics

When you have `BaseMetricCollector` from pkg `metrics`
you can call create records in Datadog.

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
[Inventory](https://github.com/coopnorge/go-logger/blob/main/docs/index.md#hooks),
[GitHub](https://inventory.internal.coop/docs/default/component/go-logger/#hooks))
to capture the `trace_id` `span_id`.

```go
package main

import (
	"github.com/coopnorge/go-datadog-lib/v2/tracelogger"
	"github.com/coopnorge/go-logger"

	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
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
	err = b(ctx)
	span.Finish(tracer.WithError(err))
}

func b(ctx context.Context) {
	logger.WithContext(ctx).Info("Hello")
}
```
