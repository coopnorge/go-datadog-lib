package grpc

import (
	"google.golang.org/grpc"
	ddGrpc "gopkg.in/DataDog/dd-trace-go.v1/contrib/google.golang.org/grpc"
)

// TraceUnaryServerInterceptor for Datadog Log Integration, middleware will create span that can be used from context
// Deprecated: Use UnaryServerInterceptor instead. This function will be removed in a later version.
func TraceUnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return ddGrpc.UnaryServerInterceptor()
}

// TraceStreamServerInterceptor for Datadog Log Integration, middleware will create span that can be used from context
// Deprecated: Use StreamServerInterceptor instead. This function will be removed in a later version.
func TraceStreamServerInterceptor() grpc.StreamServerInterceptor {
	return ddGrpc.StreamServerInterceptor()
}

// UnaryServerInterceptor returns a middleware that creates datadog-spans on incoming requests, and stores them in the requests' context.
func UnaryServerInterceptor(options ...Option) grpc.UnaryServerInterceptor {
	opts := convertOptions(options...)
	return ddGrpc.UnaryServerInterceptor(opts...)
}

// StreamServerInterceptor returns a middleware that creates datadog-spans on incoming requests, and stores them in the requests' context.
func StreamServerInterceptor(options ...Option) grpc.StreamServerInterceptor {
	opts := convertOptions(options...)
	return ddGrpc.StreamServerInterceptor(opts...)
}
