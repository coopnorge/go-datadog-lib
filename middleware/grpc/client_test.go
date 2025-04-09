package grpc_test

import (
	"context"
	"net"
	"strconv"
	"strings"
	"testing"

	datadogMiddleware "github.com/coopnorge/go-datadog-lib/v2/middleware/grpc"

	"github.com/coopnorge/go-datadog-lib/v2/internal/testhelpers"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	testgrpc "google.golang.org/grpc/interop/grpc_testing"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/test/bufconn"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/mocktracer"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

type testServer struct {
	testgrpc.UnimplementedTestServiceServer

	traceparent string
	tracestate  string
	ddTraceID   uint64
	ddParentID  uint64
}

const ddService = "DD_SERVICE"

func (s *testServer) hydrateTraceData(ctx context.Context) error {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return status.Error(codes.InvalidArgument, "not metadata in request")
	}
	s.traceparent = strings.Join(md.Get("traceparent"), "")
	s.tracestate = strings.Join(md.Get("tracestate"), "")
	ddParentID, _ := strconv.ParseUint(strings.Join(md.Get("x-datadog-parent-id"), ""), 10, 64)
	s.ddParentID = ddParentID
	ddTraceID, _ := strconv.ParseUint(strings.Join(md.Get("x-datadog-trace-id"), ""), 10, 64)
	s.ddTraceID = ddTraceID
	return nil
}

func (s *testServer) EmptyCall(ctx context.Context, _ *testgrpc.Empty) (*testgrpc.Empty, error) {
	if err := s.hydrateTraceData(ctx); err != nil {
		return nil, err
	}
	return new(testgrpc.Empty), nil
}

func (s *testServer) StreamingOutputCall(_ *testgrpc.StreamingOutputCallRequest, streamingServer grpc.ServerStreamingServer[testgrpc.StreamingOutputCallResponse]) error {
	if err := s.hydrateTraceData(streamingServer.Context()); err != nil {
		return err
	}
	return nil
}

const bufSize = 1024 * 1024

func TestTraceUnaryClientInterceptor(t *testing.T) {
	testhelpers.ConfigureDatadog(t)

	// Start Datadog tracer, so that we don't create NoopSpans.
	testTracer := mocktracer.Start()

	server := &testServer{}

	listener := bufconn.Listen(bufSize)
	grpcServer := grpc.NewServer()
	testgrpc.RegisterTestServiceServer(grpcServer, server)
	errCh := make(chan error, 1)
	go func() {
		errCh <- grpcServer.Serve(listener)
	}()

	conn, err := grpc.NewClient("dns:///localhost",
		grpc.WithContextDialer(func(_ context.Context, _ string) (net.Conn, error) { return listener.Dial() }),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithUnaryInterceptor(datadogMiddleware.TraceUnaryClientInterceptor()),
	)
	require.NoError(t, err)

	client := testgrpc.NewTestServiceClient(conn)
	span, spanCtx := tracer.StartSpanFromContext(context.Background(), ddService, tracer.ResourceName("/helloworld"))
	defer span.Finish()
	_, err = client.EmptyCall(spanCtx, &testgrpc.Empty{})
	require.NoError(t, err)

	testTracer.Stop()

	spans := testTracer.FinishedSpans()
	require.Equal(t, 1, len(spans))
	finishedSpan := spans[0]
	assert.Equal(t, finishedSpan.TraceID(), server.ddTraceID)
	assert.Equal(t, finishedSpan.SpanID(), server.ddParentID)

	assert.Empty(t, server.traceparent, "Datadog's mocktracer does not propagate W3C-headers as of writing this test. If they start propagating it, we should remove the separate test below, and update this test to assert the correct W3C-header.")
}

func TestTraceUnaryClientInterceptorW3C(t *testing.T) {
	testhelpers.ConfigureDatadog(t)

	// Start Datadog tracer, so that we don't create NoopSpans.
	// Start real tracer (not mocktracer), to propagate Traceparent.
	tracer.Start()

	server := &testServer{}
	listener := bufconn.Listen(bufSize)
	grpcServer := grpc.NewServer()
	testgrpc.RegisterTestServiceServer(grpcServer, server)
	errCh := make(chan error, 1)
	go func() {
		errCh <- grpcServer.Serve(listener)
	}()

	conn, err := grpc.NewClient("dns:///localhost",
		grpc.WithContextDialer(func(_ context.Context, _ string) (net.Conn, error) { return listener.Dial() }),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithUnaryInterceptor(datadogMiddleware.TraceUnaryClientInterceptor()),
	)
	require.NoError(t, err)

	client := testgrpc.NewTestServiceClient(conn)
	span, spanCtx := tracer.StartSpanFromContext(context.Background(), ddService, tracer.ResourceName("/helloworld"))
	defer span.Finish()
	_, err = client.EmptyCall(spanCtx, &testgrpc.Empty{})
	require.NoError(t, err)

	// Assert TraceParent
	require.NotEmpty(t, server.traceparent)
	parts := strings.Split(server.traceparent, "-")
	require.Equal(t, 4, len(parts))
	// version
	assert.Equal(t, "00", parts[0], "w3c version is not correct")
	// trace-id
	assert.Equal(t, 32, len(parts[1]), "w3c trace-id has invalid length")
	assert.NotEqual(t, "00000000000000000000000000000000", parts[1], "w3c trace-id is zero")
	// parent-id
	assert.Equal(t, 16, len(parts[2]), "w3c parent-id has invalid length")
	assert.NotEqual(t, "0000000000000000", parts[2], "w3c parent-id is zero")
	// trace-flags
	assert.Equal(t, "01", parts[3], "w3c trace-flags not is not correct")

	// Assert TraceState
	parts = strings.Split(server.tracestate, ",")
	require.True(t, len(parts) >= 1)
	found := false
	for _, listMember := range parts {
		if strings.HasPrefix(listMember, "dd=") {
			assert.NotEmpty(t, strings.TrimPrefix(listMember, "dd="))
			found = true
		}
	}
	assert.True(t, found, "Did not find Datadog's list-member in w3c tracestate")
}

func TestStreamClientInterceptor(t *testing.T) {
	testhelpers.ConfigureDatadog(t)

	// Start Datadog tracer, so that we don't create NoopSpans.
	testTracer := mocktracer.Start()

	server := &testServer{}

	listener := bufconn.Listen(bufSize)
	grpcServer := grpc.NewServer()
	testgrpc.RegisterTestServiceServer(grpcServer, server)
	errCh := make(chan error, 1)
	go func() {
		errCh <- grpcServer.Serve(listener)
	}()

	conn, err := grpc.NewClient("dns:///localhost",
		grpc.WithContextDialer(func(_ context.Context, _ string) (net.Conn, error) { return listener.Dial() }),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithStreamInterceptor(datadogMiddleware.StreamClientInterceptor()),
	)
	require.NoError(t, err)

	client := testgrpc.NewTestServiceClient(conn)
	span, spanCtx := tracer.StartSpanFromContext(context.Background(), ddService, tracer.ResourceName("/helloworld"))
	c, err := client.StreamingOutputCall(spanCtx, &testgrpc.StreamingOutputCallRequest{})
	require.NoError(t, err)
	span.Finish()

	c.Recv()

	testTracer.Stop()

	spans := testTracer.FinishedSpans()
	require.Equal(t, 4, len(spans))
	for _, finishedSpan := range spans {
		assert.Equal(t, finishedSpan.TraceID(), server.ddTraceID)
		if finishedSpan.OperationName() == ddService {
			assert.Equal(t, finishedSpan.ParentID(), uint64(0))
		} else if finishedSpan.OperationName() == "grpc.client" {
			assert.Equal(t, finishedSpan.ParentID(), finishedSpan.TraceID())
			assert.Equal(t, finishedSpan.ParentID(), server.ddTraceID)
		} else {
			assert.Equal(t, "grpc.message", finishedSpan.OperationName())
			assert.Equal(t, finishedSpan.ParentID(), server.ddParentID)
		}
	}
}
