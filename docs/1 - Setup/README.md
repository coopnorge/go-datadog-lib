# Setup

Prepare configuration for your container.

## 1. Setup container

Following configuration example related to Kubernetes and Kustomize.

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: my-app
  labels:
    app: my-app
    version: v1
    tags.datadoghq.com/service: "my-app-api"
    tags.datadoghq.com/version: "v1"
spec:
  selector:
    matchLabels:
      app: my-app
      version: v1
  template:
    metadata:
      labels:
        app: my-app
        version: v1
        tags.datadoghq.com/service: "my-app-api"
        tags.datadoghq.com/version: "v1"
      annotations:
        proxy.istio.io/config: '{ "holdApplicationUntilProxyStarts": true }'
    spec:
      serviceAccountName: my-app-app
      containers:
        - name: my-app-api
          envFrom:
            - secretRef:
                name: external-secrets-my-app
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

Most important is that your paths for sockets are the same as this.
The reason for that is Datadog implementation in go would not connect
to APM if you will pass Unix socket prefix.

```yaml
- name: DD_DOGSTATSD_URL
  value: "unix:///var/run/datadog/dsd.socket"
- name: DD_TRACE_AGENT_URL
  value: "/var/run/datadog/apm.socket"
```

It's how the application will be shown in Datadog APM.

```yaml
- name: DD_SERVICE
  valueFrom:
    fieldRef:
      fieldPath: metadata.labels['tags.datadoghq.com/service']
```

Depending on your policy, you can have an API version or tag/commit from git.

```yaml
- name: DD_VERSION
  valueFrom:
    fieldRef:
      fieldPath: metadata.labels['tags.datadoghq.com/version']
```

## 2. Application setup

Create pkg configuration for boostrap datadog.

```go
package main

import (
	"github.com/coopnorge/go-datadog-lib"
	"github.com/coopnorge/go-datadog-lib/config"
)

func main() {
	// Your app initialization
	/// ... 

	// From your core configuration add datadog related values
	ddCfg := config.DatadogConfig{
		Env:            "dd_env",
		Service:        "dd_service",
		ServiceVersion: "dd_version",
		DSD:            "dd_dogstatsd_url",
		APM:            "dd_trace_agent_url",
	}

	// When you start other processes start datadog
	withExtraProfiler := true
	go_datadog_lib.StartDatadog(ddCfg, withExtraProfiler)
	
	// Stop datadog with yours other processes
	handleGracefulShutdown(go_datadog_lib.GracefulDatadogShutdown)
}
```

## 3. Middleware

To have better tracing you need add to your gRPC custom
middleware that will extend context.

It's needed to relate logs with your trace data in APM.

To do that simple add Go - Datadog middleware to your gRPC interceptor.

Take a look `"github.com/coopnorge/go-datadog-lib/grpc"`
function `TraceUnaryServerInterceptor`

```go
    // This is gRPC server configuration builder
	cfgBuilder.AddGrpcUnaryInterceptors(
		grpctrace.UnaryServerInterceptor(
			grpctrace.WithServiceName(cfg.DatadogService),
			grpctrace.WithStreamCalls(false),
		),
		grpc.TraceUnaryServerInterceptor(),
	)

    // This middleware will extend context for tracing and logs
    // grpc.TraceUnaryServerInterceptor()
```

## Common issue

After that Datadog will try to connect to the socket
and will start to send all information in the background.

In different setup, you could have error logs
that Datadog cannot connect to the socket and
tried to connect via HTTP. That could be related
to issue when your container starts faster
and sockets were not ready to communicate with Agent
or Agent was started later.

