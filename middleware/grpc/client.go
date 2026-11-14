// Package grpc provides Datadog tracing interceptors for gRPC servers and
// clients.
//
// Deprecated: The whole module is deprecated. Use otelgrpc from
// go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc:
// grpc.StatsHandler(otelgrpc.NewServerHandler()) for servers and
// grpc.WithStatsHandler(otelgrpc.NewClientHandler()) for clients. See the root
// package documentation for the migration guide and ADR.
package grpc

import (
	"context"

	"github.com/coopnorge/go-datadog-lib/v2/internal"
	"google.golang.org/grpc"

	ddGrpc "github.com/DataDog/dd-trace-go/contrib/google.golang.org/grpc/v2"
)

// TraceUnaryClientInterceptor create a client-interceptor to automatically create child-spans, and append to gRPC metadata.
//
// Deprecated: Use UnaryClientInterceptor instead. This function will be removed in a later version.
// The whole module is deprecated; use otelgrpc.NewClientHandler() from
// go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc, passed to
// grpc.WithStatsHandler.
func TraceUnaryClientInterceptor() grpc.UnaryClientInterceptor {
	if internal.IsDatadogDisabled() {
		return noOpUnaryClientInterceptor()
	}

	return ddGrpc.UnaryClientInterceptor()
}

// UnaryClientInterceptor create a client-interceptor to automatically create child-spans, and append to gRPC metadata.
// UnaryServerInterceptor returns a middleware that creates datadog-spans on outgoing requests, and adds them to the request's gRPC-metadata.
//
// Deprecated: The whole module is deprecated. Use otelgrpc.NewClientHandler() from
// go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc, passed to
// grpc.WithStatsHandler.
func UnaryClientInterceptor(options ...Option) grpc.UnaryClientInterceptor {
	if internal.IsDatadogDisabled() {
		return noOpUnaryClientInterceptor()
	}

	opts := convertOptions(options...)
	return ddGrpc.UnaryClientInterceptor(opts...)
}

// StreamClientInterceptor create a client-interceptor to automatically create child-spans, and append to gRPC metadata.
//
// Deprecated: The whole module is deprecated. Use otelgrpc.NewClientHandler() from
// go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc, passed to
// grpc.WithStatsHandler.
func StreamClientInterceptor(options ...Option) grpc.StreamClientInterceptor {
	if internal.IsDatadogDisabled() {
		return noOpStreamClientInterceptor()
	}
	opts := convertOptions(options...)
	return ddGrpc.StreamClientInterceptor(opts...)
}

func noOpStreamClientInterceptor() grpc.StreamClientInterceptor {
	return func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
		return streamer(ctx, desc, cc, method, opts...)
	}
}

func noOpUnaryClientInterceptor() grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply any, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		return invoker(ctx, method, req, reply, cc, opts...)
	}
}
