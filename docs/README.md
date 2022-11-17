# Datadog is an observability service for cloud-scale applications

## Go Datadog lib

> Reduces the complexity of initializing
> these services inside your application.
> Also provides abstract code to work with metrics.

## APM

In Coop No our default setup is tracing
applications with CPU profiling support,
that is enabled by default in the package bootstrap.

## Custom metrics - DD StatsD

Inside the `metric` package you can find
the base client for Datadog StatsD and
simple metric service
that allows sending `Incr`, `Gauge`, and `Count`.

## How Datadog works

![diagram_dd_com](dd_com_app.png)
