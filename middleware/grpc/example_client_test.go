package grpc_test

import (
	"context"

	datadogMiddleware "github.com/coopnorge/go-datadog-lib/v2/middleware/grpc"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/testing/testpb"
	"google.golang.org/grpc"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

func ExampleUnaryClientInterceptor() {
	ctx := context.Background()

	cc, err := grpc.NewClient(
		"https://example.com",
		grpc.WithUnaryInterceptor(datadogMiddleware.UnaryClientInterceptor()),
	)
	if err != nil {
		panic(err)
	}
	defer cc.Close()

	client := testpb.NewTestServiceClient(cc)

	span, ctx := tracer.StartSpanFromContext(ctx, "grpc.request")
	resp, err := client.Ping(ctx, &testpb.PingRequest{})
	span.Finish(tracer.WithError(err))
	if err != nil {
		panic(err)
	}
	println(resp)
}

func ExampleStreamClientInterceptor() {
	ctx := context.Background()

	cc, err := grpc.NewClient(
		"https://example.com",
		grpc.WithStreamInterceptor(datadogMiddleware.StreamClientInterceptor()),
	)
	if err != nil {
		panic(err)
	}
	defer cc.Close()

	client := testpb.NewTestServiceClient(cc)

	span, ctx := tracer.StartSpanFromContext(ctx, "grpc.stream")
	stream, err := client.PingStream(ctx)
	defer span.Finish()
	if err != nil {
		span.Finish(tracer.WithError(err))
		panic(err)
	}
	resp, err := stream.Recv()
	if err != nil {
		span.Finish(tracer.WithError(err))
		panic(err)
	}
	println(resp)
}
