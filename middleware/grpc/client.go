package grpc

import (
	"google.golang.org/grpc"

	ddGrpc "gopkg.in/DataDog/dd-trace-go.v1/contrib/google.golang.org/grpc"
)

// TraceUnaryClientInterceptor create a client-interceptor to automatically create child-spans, and append to gRPC metadata.
func TraceUnaryClientInterceptor() grpc.UnaryClientInterceptor {
	return ddGrpc.UnaryClientInterceptor()
}
