# Coop Datadog Go package

![Test](https://github.com/coopnorge/go-datadog-lib/actions/workflows/test.yml/badge.svg)
![Build](https://github.com/coopnorge/go-datadog-lib/actions/workflows/build.yml/badge.svg)

Plug and play package that wraps base functionally and initialization of
Datadog Service.

- APM, StatsD Initialization
- StatsD metrics unification

Supported middleware to correlate/extend traceability and logs in Datadog.

- [X] gRPC Unary Server
- [X] HTTP - Echo

## Documentation

There is detailed documentation stored in [docs](docs/).

## Mocks

To generate or update mocks use tools
[Eitri](https://github.com/Clink-n-Clank/Eitri) or use directly
[Mockhandler](github.com/sanposhiho/gomockhandle)

## Development workflow

### Validate

```bash
docker compose run --rm golang-devtools validate
```

### Other targets

```bash
docker compose run --rm golang-devtools help
```

## User documentation

User documentation is build using TechDocs and published to
[Inventory](https://inventory.internal.coop/docs/default/component/go-datadog-lib).

To list the commands available for the TechDocs image:

```sh
docker compose run --rm help
```

For more information see the [TechDocs Engineering
Image](https://github.com/coopnorge/engineering-docker-images/tree/main/images/techdocs).

### Documentation validation

To Validate changed documentation:

```sh
docker compose run --rm techdocs validate
```

To validate all documentation:

```sh
docker compose run --rm techdocs validate MARKDOWN_FILES=docs/
```

### Documentation preview

To preview the documentation:

```sh
docker compose up techdocs
```

## Experimental Tracing

For advanced tracing capabilities and features, you can enable experimental
tracing by setting the environment variable
"DD_EXPERIMENTAL_TRACING_ENABLED" to "true", "TRUE", or "True".

When this environment variable is set, the package utilizes an experimental
tracing interceptor from
[dd-trace-go](https://github.com/DataDog/dd-trace-go/blob/main/contrib/google.golang.org/grpc/server.go)
for gRPC servers, enhancing distributed tracing capabilities.
Additionally, for the Echo middleware, the package leverages
the interceptor from
[dd-trace-go Echo Middleware](https://github.com/DataDog/dd-trace-go/tree/main/contrib/labstack/echo.v4).

In a future release, these experimental tracing features will become
the default behavior, eliminating the need for the feature flag.
