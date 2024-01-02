package grpc

import (
	"google.golang.org/grpc"

	ddGrpc "gopkg.in/DataDog/dd-trace-go.v1/contrib/google.golang.org/grpc"
)

// TraceUnaryClientInterceptor create a client-interceptor to automatically create child-spans, and append to gRPC metadata.
// Deprecated: Use UnaryClientInterceptor instead. This function will be removed in a later version.
func TraceUnaryClientInterceptor() grpc.UnaryClientInterceptor {
	return ddGrpc.UnaryClientInterceptor()
}

// UnaryClientInterceptor create a client-interceptor to automatically create child-spans, and append to gRPC metadata.
// UnaryServerInterceptor returns a middleware that creates datadog-spans on outgoing requests, and adds them to the request's gRPC-metadata.
func UnaryClientInterceptor(options ...Option) grpc.UnaryClientInterceptor {
	opts := convertOptions(options...)
	return ddGrpc.UnaryClientInterceptor(opts...)
}
