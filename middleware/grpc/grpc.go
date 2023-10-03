package grpc

import (
	"context"
	"google.golang.org/grpc/metadata"
	"strconv"

	"github.com/coopnorge/go-datadog-lib/v2/internal"
	"github.com/coopnorge/go-datadog-lib/v2/tracing"

	"google.golang.org/grpc"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"

	grpcmw "github.com/grpc-ecosystem/go-grpc-middleware"
)

// TraceUnaryServerInterceptor for Datadog Log Integration, middleware will create span that can be used from context
func TraceUnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(reqCtx context.Context, req interface{}, info *grpc.UnaryServerInfo, grpcReqHandler grpc.UnaryHandler) (interface{}, error) {
		span, spanCtx := tracer.StartSpanFromContext(reqCtx, info.FullMethod, tracer.ResourceName("grpc.request"))
		defer span.Finish()

		extCtx := internal.ExtendedContextWithMetadata(spanCtx, internal.TraceContextKey{}, tracing.TraceDetails{DatadogSpan: span})

		md := metadata.New(map[string]string{"traceID": strconv.FormatUint(span.Context().TraceID(), 10)})
		grpc.SetHeader(extCtx, md)

		return grpcReqHandler(extCtx, req)
	}
}

// TraceStreamServerInterceptor for Datadog Log Integration, middleware will create span that can be used from context
func TraceStreamServerInterceptor() grpc.StreamServerInterceptor {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		span, spanCtx := tracer.StartSpanFromContext(ss.Context(), info.FullMethod, tracer.ResourceName("grpc.request"))
		defer span.Finish()

		extCtx := internal.ExtendedContextWithMetadata(spanCtx, internal.TraceContextKey{}, tracing.TraceDetails{DatadogSpan: span})

		return handler(srv, &grpcmw.WrappedServerStream{
			ServerStream:   ss,
			WrappedContext: extCtx,
		})
	}
}
