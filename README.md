# Coop Datadog Go package

![Build](https://github.com/coopnorge/go-datadog-lib/actions/workflows/cicd.yaml/badge.svg)

Plug and play package that wraps base functionally and initialization of Datadog
Service.

- APM, StatsD Initialization
- StatsD metrics unification

Supported middleware to correlate/extend traceability and logs in Datadog.

- [x] gRPC Server
- [x] gRPC Client
- [x] HTTP - Echo
- [x] HTTP - Standard library Client
- [x] Database - GORM
- [x] Database - Standard library

## Documentation

<https://pkg.go.dev/github.com/coopnorge/go-datadog-lib/v2>

## Development workflow

The source code is build using `mage`.

### Prerequisites

1. Install Go version 1.24 or later and
   [Docker](https://docs.docker.com/get-docker/).

2. Install Go tools:

   ```console
   go install tool
   ```

### Validate

```console
go tool mage validate
```

### Other targets

```console
go tool mage -l
```

### Mocks

To generate or update mocks use
[`gomockhandler`](github.com/sanposhiho/gomockhandler) together with
[Uber's `mockgen`](go.uber.org/mock/mockge).

#### Check mocks

```console
go install tool
go tool github.com/sanposhiho/gomockhandler -config ./gomockhandler.json check
```

#### Generate / Update mocks

```console
go install tool
go tool github.com/sanposhiho/gomockhandler -config ./gomockhandler.json mockgen
```

## User documentation

User documentation is built using TechDocs and published to
[Inventory](https://inventory.internal.coop/docs/default/component/go-datadog-lib).

To list the commands available for the TechDocs image:

```console
docker compose run --rm help
```

For more information see the
[TechDocs Engineering Image](https://github.com/coopnorge/engineering-docker-images/tree/main/images/techdocs).

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
