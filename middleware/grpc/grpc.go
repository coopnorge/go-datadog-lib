package grpc

import (
	"context"

	"github.com/coopnorge/go-datadog-lib/v2/internal"
	"github.com/coopnorge/go-datadog-lib/v2/tracing"
	grpcmw "github.com/grpc-ecosystem/go-grpc-middleware"
	"google.golang.org/grpc"
	ddGrpc "gopkg.in/DataDog/dd-trace-go.v1/contrib/google.golang.org/grpc"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

// TraceUnaryServerInterceptor for Datadog Log Integration, middleware will create span that can be used from context
func TraceUnaryServerInterceptor() grpc.UnaryServerInterceptor {
	if internal.IsExperimentalTracingEnabled() {
		return traceUnaryServerInterceptorExperimental()
	}
	return traceUnaryServerInterceptorLegacy()
}

// TraceStreamServerInterceptor for Datadog Log Integration, middleware will create span that can be used from context
func TraceStreamServerInterceptor() grpc.StreamServerInterceptor {
	if internal.IsExperimentalTracingEnabled() {
		return traceStreamServerInterceptorExperimental()
	}
	return traceStreamServerInterceptorLegacy()
}

func traceUnaryServerInterceptorLegacy() grpc.UnaryServerInterceptor {
	return func(reqCtx context.Context, req interface{}, info *grpc.UnaryServerInfo, grpcReqHandler grpc.UnaryHandler) (interface{}, error) {
		span, spanCtx := tracer.StartSpanFromContext(reqCtx, info.FullMethod, tracer.ResourceName("grpc.request"))
		defer span.Finish()

		extCtx := internal.ExtendedContextWithMetadata(spanCtx, internal.TraceContextKey{}, tracing.TraceDetails{DatadogSpan: span})

		return grpcReqHandler(extCtx, req)
	}
}

func traceUnaryServerInterceptorExperimental() grpc.UnaryServerInterceptor {
	return ddGrpc.UnaryServerInterceptor()
}

func traceStreamServerInterceptorLegacy() grpc.StreamServerInterceptor {
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

func traceStreamServerInterceptorExperimental() grpc.StreamServerInterceptor {
	return ddGrpc.StreamServerInterceptor()
}
