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
func TraceUnaryClientInterceptor() grpc.UnaryClientInterceptor {
	if internal.IsDatadogDisabled() {
		return noOpUnaryClientInterceptor()
	}

	return ddGrpc.UnaryClientInterceptor()
}

// UnaryClientInterceptor create a client-interceptor to automatically create child-spans, and append to gRPC metadata.
// UnaryServerInterceptor returns a middleware that creates datadog-spans on outgoing requests, and adds them to the request's gRPC-metadata.
func UnaryClientInterceptor(options ...Option) grpc.UnaryClientInterceptor {
	if internal.IsDatadogDisabled() {
		return noOpUnaryClientInterceptor()
	}

	opts := convertOptions(options...)
	return ddGrpc.UnaryClientInterceptor(opts...)
}

// StreamClientInterceptor create a client-interceptor to automatically create child-spans, and append to gRPC metadata.
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
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		return invoker(ctx, method, req, reply, cc, opts...)
	}
}
