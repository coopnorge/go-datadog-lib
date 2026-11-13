# Migrate to OpenTelemetry

> [!WARNING]
>
> This library is deprecated. See the [go-telemetry-lib ADR 0001][adr-0001] for
> the reasoning and the migration steps below for replacements.
>
> The repository will be archived on 2027-03-01. Until then it only receives
> non-breaking dependency updates.

## Why

Coop prefers official OpenTelemetry tooling over internal wrappers. Services use
the OpenTelemetry API and instrumentation packages directly, and
`github.com/coopnorge/go-telemetry-lib` handles startup and provides a log hook.
It is not a 1:1 replacement for `github.com/coopnorge/go-datadog-lib/v2`.

- [go-telemetry-lib ADR 0001][adr-0001]: deprecate `go-datadog-lib`.
- [cloud-platform ADR 0007][adr-0007]: OpenTelemetry observability SDK.
- [go-telemetry-lib getting started][getting-started].

## Timeline and support

- No new features and no bug fixes.
- Non-breaking dependency updates only, until 2027-03-01.
- The repository is archived on 2027-03-01.
- Services that need continuous profiling stay on this library until then.

## Before you start

`telemetry.Start` reads its configuration from environment variables:

- `OTEL_SERVICE_NAME`: service name.
- `OTEL_RESOURCE_ATTRIBUTES`: must contain `deployment.environment.name` and
  `service.version`.
- `OTEL_EXPORTER_OTLP_ENDPOINT`: where telemetry is exported.
- `OTEL_SDK_DISABLED`: disables the SDK when set to `true`.

`telemetry.Start` returns an error if the service name, the environment or the
version is missing.

The Helm charts `coop-app-chart` v0.40.0 or later and `coop-job-chart` v0.12.0
or later inject these variables. With older chart versions, set them yourself.

Create the shutdown context with `signals.HandleSignals` from
`github.com/coopnorge/go-signal-lib` and pass it to `telemetry.Start`.

## Mapping

<!-- markdownlint-disable MD013 -->

| `go-datadog-lib`                                                 | Replacement                                                                                                                                                                                                                |
| ---------------------------------------------------------------- | -------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `Start`, `StartDatadog`, `StopFunc`                              | `telemetry.Start(ctx)` from `github.com/coopnorge/go-telemetry-lib`. It returns `StopFunc func() error`, which takes no context.                                                                                           |
| `config`                                                         | No config struct. `telemetry.Start` reads the `OTEL_*` environment variables.                                                                                                                                              |
| `tracelogger`                                                    | `loghook.NewHook()` from `github.com/coopnorge/go-telemetry-lib/loghook`, registered with `logger.ConfigureGlobalLogger(logger.WithHook(loghook.NewHook()))`.                                                              |
| `tracing.CreateChildSpan`, `tracing.CreateNestedTrace`           | `otel.Tracer("name").Start(ctx, "span")` from `go.opentelemetry.io/otel`. Nesting follows the context.                                                                                                                     |
| `tracing.LogWithTrace`, `tracing.LogFieldsWithTrace`             | The log hook adds `trace_id` and `span_id` to log entries.                                                                                                                                                                 |
| `errors.ErrorHandler`                                            | `span.RecordError(err)` and `span.SetStatus(codes.Error, msg)` (`go.opentelemetry.io/otel/codes`).                                                                                                                         |
| `metric`, `metrics`                                              | Instruments from `otel.Meter("name")` (`go.opentelemetry.io/otel/metric`). Create them once, keep them and reuse them. No wrapper is offered, on purpose.                                                                  |
| `middleware/grpc`                                                | `go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc`. Use `grpc.NewServer(grpc.StatsHandler(otelgrpc.NewServerHandler()))`. For clients use `grpc.WithStatsHandler(otelgrpc.NewClientHandler())`. |
| `middleware/http`                                                | `go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp`. Use `otelhttp.NewHandler(mux, "http-server")`. For clients use `otelhttp.NewTransport(http.DefaultTransport)`.                                            |
| `middleware/echo` (Echo v4)                                      | `github.com/labstack/echo-otel/v4`. Use `e.Use(echootel.NewMiddleware("app.example.com"))`.                                                                                                                                |
| `middleware/gorm`, and `middleware/database` when used with GORM | `gorm.io/plugin/opentelemetry/tracing`. Use `db.Use(tracing.NewPlugin())` after `gorm.Open`.                                                                                                                               |

<!-- markdownlint-enable MD013 -->

Notes:

- The `serverName` argument to `echootel.NewMiddleware` is a host with an
  optional port. Spans are named after the matched route, for example
  `GET /users/:id`, and carry `http.route`.
- `github.com/labstack/echo-otel` follows its own versioning, which is not
  SemVer. A MINOR release may include breaking API changes, so pin the version.
- With GORM, call `db.WithContext(ctx)` so spans attach to the request trace.

## Examples

### Startup, shutdown and log hook

Before:

```go
import (
	"context"

	coopdatadog "github.com/coopnorge/go-datadog-lib/v2"
)

func run() error {
	stop, err := coopdatadog.Start(context.Background())
	if err != nil {
		return err
	}
	defer func() {
		if err := stop(); err != nil {
			panic(err)
		}
	}()

	// ...

	return nil
}
```

After:

```go
import (
	"context"

	"github.com/coopnorge/go-logger"
	signals "github.com/coopnorge/go-signal-lib"
	telemetry "github.com/coopnorge/go-telemetry-lib"
	"github.com/coopnorge/go-telemetry-lib/loghook"
)

func main() {
	ctx, cancel := signals.HandleSignals(context.Background())
	defer cancel()

	logger.ConfigureGlobalLogger(logger.WithHook(loghook.NewHook()))

	stop, err := telemetry.Start(ctx)
	if err != nil {
		logger.Fatal(err)
	}
	defer func() {
		if err := stop(); err != nil {
			logger.Fatal(err)
		}
	}()

	// ...
}
```

### Manual spans and errors

Before:

```go
import "github.com/coopnorge/go-datadog-lib/v2/tracing"

func doWork(ctx context.Context) error {
	span := tracing.CreateChildSpan(ctx, "work", "doWork")
	defer span.Finish()

	// ...
}
```

After:

```go
import (
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
)

var tracer = otel.Tracer("my-service/work")

func doWork(ctx context.Context) error {
	ctx, span := tracer.Start(ctx, "doWork")
	defer span.End()

	if err := step(ctx); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return err
	}
	return nil
}
```

### Metrics

Before:

```go
import "github.com/coopnorge/go-datadog-lib/v2/metrics"

metrics.Count("orders.created", 1)
metrics.Histogram("orders.size", size)
```

After:

```go
import (
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/metric"
)

meter := otel.Meter("my-service/orders")

ordersCreated, err := meter.Int64Counter("orders.created")
if err != nil {
	return err
}
orderSize, err := meter.Float64Histogram("orders.size")
if err != nil {
	return err
}

ordersCreated.Add(ctx, 1)
orderSize.Record(ctx, size)
```

Create the instruments once, for example at startup, and reuse them.

### gRPC

Before:

```go
import datadogMiddleware "github.com/coopnorge/go-datadog-lib/v2/middleware/grpc"

server := grpc.NewServer(
	grpc.UnaryInterceptor(datadogMiddleware.UnaryServerInterceptor()),
	grpc.StreamInterceptor(datadogMiddleware.StreamServerInterceptor()),
)

conn, err := grpc.NewClient(target,
	grpc.WithUnaryInterceptor(datadogMiddleware.UnaryClientInterceptor()),
	grpc.WithStreamInterceptor(datadogMiddleware.StreamClientInterceptor()),
)
```

After:

```go
import "go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"

server := grpc.NewServer(grpc.StatsHandler(otelgrpc.NewServerHandler()))

conn, err := grpc.NewClient(target,
	grpc.WithStatsHandler(otelgrpc.NewClientHandler()),
)
```

### HTTP

Before:

```go
import datadogMiddleware "github.com/coopnorge/go-datadog-lib/v2/middleware/http"

client := datadogMiddleware.AddTracingToClient(&http.Client{})
```

After:

```go
import "go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"

handler := otelhttp.NewHandler(mux, "http-server")

client := &http.Client{
	Transport: otelhttp.NewTransport(http.DefaultTransport),
}
```

### Echo v4

Before:

```go
import coopEchoDatadog "github.com/coopnorge/go-datadog-lib/v2/middleware/echo"

e := echo.New()
e.Use(coopEchoDatadog.TraceServerMiddleware())
```

After:

```go
import echootel "github.com/labstack/echo-otel/v4"

e := echo.New()
e.Use(echootel.NewMiddleware("app.example.com"))
```

### GORM

Before:

```go
import ddGorm "github.com/coopnorge/go-datadog-lib/v2/middleware/gorm"

db, err := ddGorm.NewORM(mysql.Open(dsn), &gorm.Config{})
```

After:

```go
import "gorm.io/plugin/opentelemetry/tracing"

db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
if err != nil {
	return err
}
if err := db.Use(tracing.NewPlugin()); err != nil {
	return err
}

db.WithContext(ctx).First(&user)
```

## No direct replacement

These features have no equivalent in the OpenTelemetry packages listed above:

- `WithNonErrorCodes` in `middleware/grpc`.
- `metrics.SimpleEvent`.
- `tracing.AppendUserToTrace`.
- `tracing.OverrideTraceResourceName`.
- Tracing of plain `database/sql` without GORM.
- Continuous profiling. Stay on this library until it is archived.

## Suggested migration order

1. Set the `OTEL_*` environment variables, or upgrade the Helm chart.
2. Replace `coopdatadog.Start` with `telemetry.Start`.
3. Add the log hook.
4. Replace the middleware.
5. Replace manual spans and error handling.
6. Replace metrics.
7. Remove the `github.com/coopnorge/go-datadog-lib/v2` dependency.

## Getting help

See the [go-telemetry-lib getting started page][getting-started].

[adr-0001]:
  https://inventory.internal.coop/catalog/default/component/go-telemetry-lib/adrs?record=0001-deprecate-go-datadog-lib.md
[adr-0007]:
  https://inventory.internal.coop/catalog/default/system/cloud-platform/adrs?record=0007-opentelemetry-observability-sdk.md
[getting-started]:
  https://inventory.internal.coop/docs/default/component/go-telemetry-lib/getting-started/
