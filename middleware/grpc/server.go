package grpc

import (
	"context"

	"github.com/coopnorge/go-datadog-lib/v2/internal"
	"google.golang.org/grpc"
	ddGrpc "gopkg.in/DataDog/dd-trace-go.v1/contrib/google.golang.org/grpc"
)

// TraceUnaryServerInterceptor for Datadog Log Integration, middleware will create span that can be used from context
//
// Deprecated: Use UnaryServerInterceptor instead. This function will be removed in a later version.
func TraceUnaryServerInterceptor() grpc.UnaryServerInterceptor {
	if internal.IsDatadogDisabled() {
		return noOpUnaryServerInterceptor()
	}

	return ddGrpc.UnaryServerInterceptor()
}

// TraceStreamServerInterceptor for Datadog Log Integration, middleware will create span that can be used from context
//
// Deprecated: Use StreamServerInterceptor instead. This function will be removed in a later version.
func TraceStreamServerInterceptor() grpc.StreamServerInterceptor {
	if internal.IsDatadogDisabled() {
		return noOpStreamServerInterceptor()
	}

	return ddGrpc.StreamServerInterceptor()
}

// UnaryServerInterceptor returns a middleware that creates datadog-spans on incoming requests, and stores them in the requests' context.
func UnaryServerInterceptor(options ...Option) grpc.UnaryServerInterceptor {
	if internal.IsDatadogDisabled() {
		return noOpUnaryServerInterceptor()
	}
	opts := convertOptions(options...)
	return ddGrpc.UnaryServerInterceptor(opts...)
}

// StreamServerInterceptor returns a middleware that creates datadog-spans on incoming requests, and stores them in the requests' context.
func StreamServerInterceptor(options ...Option) grpc.StreamServerInterceptor {
	if internal.IsDatadogDisabled() {
		return noOpStreamServerInterceptor()
	}
	opts := convertOptions(options...)
	return ddGrpc.StreamServerInterceptor(opts...)
}

func noOpUnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, _ *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		return handler(ctx, req)
	}
}

func noOpStreamServerInterceptor() grpc.StreamServerInterceptor {
	return func(srv interface{}, ss grpc.ServerStream, _ *grpc.StreamServerInfo, handler grpc.StreamHandler) (err error) {
		return handler(srv, ss)
	}
}
