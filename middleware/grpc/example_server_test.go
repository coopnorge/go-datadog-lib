package grpc_test

import (
	"context"
	"fmt"
	"net"

	coopdatadog "github.com/coopnorge/go-datadog-lib/v2"
	datadogMiddleware "github.com/coopnorge/go-datadog-lib/v2/middleware/grpc"
	"google.golang.org/grpc"
)

func Example() {
	err := run()
	if err != nil {
		panic(err)
	}
}

func run() error {
	ctx := context.Background()
	stop, err := coopdatadog.Start(ctx)
	if err != nil {
		panic(err)
	}
	defer func() {
		err := stop()
		if err != nil {
			panic(err)
		}
	}()

	ddOpts := []datadogMiddleware.Option{
		// ...
	}
	serverOpts := []grpc.ServerOption{
		grpc.UnaryInterceptor(datadogMiddleware.UnaryServerInterceptor(ddOpts...)),
		grpc.StreamInterceptor(datadogMiddleware.StreamServerInterceptor(ddOpts...)),
	}

	grpcServer := grpc.NewServer(serverOpts...)

	listener, err := (&net.ListenConfig{}).Listen(ctx, "tcp", "127.0.0.1")
	if err != nil {
		return fmt.Errorf("failed to start tcp listener: %w", err)
	}

	err = grpcServer.Serve(listener)
	if err != nil {
		return nil
	}

	return nil
}
