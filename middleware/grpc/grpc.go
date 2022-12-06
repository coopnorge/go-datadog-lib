package grpc

import (
    "context"

    "github.com/coopnorge/go-datadog-lib/internal"
    "github.com/coopnorge/go-datadog-lib/tracing"

    "google.golang.org/grpc"
    "gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

// TraceUnaryServerInterceptor for Datadog Log Integration, middleware will create span that can be used from context
func TraceUnaryServerInterceptor() grpc.UnaryServerInterceptor {
    return func(reqCtx context.Context, req interface{}, info *grpc.UnaryServerInfo, grpcReqHandler grpc.UnaryHandler) (interface{}, error) {
        span, spanCtx := tracer.StartSpanFromContext(reqCtx, info.FullMethod, tracer.ResourceName("grpc.request"))
        defer span.Finish()

        extCtx := internal.ExtendedContextWithMetadata(spanCtx, internal.TraceContextKey{}, tracing.TraceDetails{DatadogSpan: span})

        return grpcReqHandler(extCtx, req)
    }
}
