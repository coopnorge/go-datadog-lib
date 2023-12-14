package grpc

import (
	"google.golang.org/grpc"
	ddGrpc "gopkg.in/DataDog/dd-trace-go.v1/contrib/google.golang.org/grpc"
)

// TraceUnaryServerInterceptor for Datadog Log Integration, middleware will create span that can be used from context
func TraceUnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return ddGrpc.UnaryServerInterceptor()
}

// TraceStreamServerInterceptor for Datadog Log Integration, middleware will create span that can be used from context
func TraceStreamServerInterceptor() grpc.StreamServerInterceptor {
	return ddGrpc.StreamServerInterceptor()
}
