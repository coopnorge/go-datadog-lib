# Coop Datadog Go package

![Test](https://github.com/coopnorge/go-datadog-lib/actions/workflows/test.yml/badge.svg)
![Build](https://github.com/coopnorge/go-datadog-lib/actions/workflows/build.yml/badge.svg)

Plug and play package that wraps base functionally
and initialization of Datadog Service.

- APM, StatsD Initialization
- StatsD metrics unification

Supported middleware to correlate/extend
traceability and logs in Datadog.

- [X] gRPC Unary Server
- [X] HTTP - Echo

## Mocks

To generate or update mocks use tools
[Eitri](https://github.com/Clink-n-Clank/Eitri)
or use directly
[Mockhandler](github.com/sanposhiho/gomockhandle)
