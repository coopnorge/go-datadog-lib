package grpc

import (
	"context"

	"github.com/coopnorge/go-datadog-lib/v2/internal"
	"google.golang.org/grpc"

	ddGrpc "gopkg.in/DataDog/dd-trace-go.v1/contrib/google.golang.org/grpc"
)

// TraceUnaryClientInterceptor create a client-interceptor to automatically create child-spans, and append to gRPC metadata.
func TraceUnaryClientInterceptor() grpc.UnaryClientInterceptor {
	if !internal.IsExperimentalTracingEnabled() {
		return noopInterceptor()
	}
	return ddGrpc.UnaryClientInterceptor()
}

func noopInterceptor() grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply any, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		return invoker(ctx, method, req, reply, cc, opts...)
	}
}
