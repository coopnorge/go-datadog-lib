package grpc

import (
	"context"
	"net"
	"sort"
	"strconv"
	"strings"
	"testing"

	"github.com/coopnorge/go-datadog-lib/v2/internal"

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
	ddTraceID   string
	ddParentID  string
}

func (s *testServer) EmptyCall(ctx context.Context, _ *testgrpc.Empty) (*testgrpc.Empty, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Error(codes.InvalidArgument, "not metadata in request")
	}
	s.traceparent = strings.Join(md.Get("Traceparent"), "")
	s.tracestate = strings.Join(md.Get("Tracestate"), "")
	s.ddTraceID = strings.Join(md.Get("X-Datadog-Trace-Id"), "")
	s.ddParentID = strings.Join(md.Get("X-Datadog-Parent-Id"), "")
	return new(testgrpc.Empty), nil
}

const bufSize = 1024 * 1024

func TestTraceUnaryClientInterceptor(t *testing.T) {
	ctx := context.Background()

	// Ensure valid datadog config, even if we don't have a datadog agent running, to fully instrument the application.
	t.Setenv("DD_ENV", "unittest")
	require.True(t, internal.IsDatadogConfigured())
	t.Setenv("DD_EXPERIMENTAL_TRACING_ENABLED", "true")
	require.True(t, internal.IsExperimentalTracingEnabled())

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

	conn, err := grpc.DialContext(ctx, "bufnet",
		grpc.WithContextDialer(func(ctx context.Context, s string) (net.Conn, error) { return listener.Dial() }),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithUnaryInterceptor(TraceUnaryClientInterceptor()),
	)
	require.NoError(t, err)

	client := testgrpc.NewTestServiceClient(conn)
	span, spanCtx := tracer.StartSpanFromContext(context.Background(), "grpc.request", tracer.ResourceName("/helloworld"))
	defer span.Finish()
	_, err = client.EmptyCall(spanCtx, &testgrpc.Empty{})
	require.NoError(t, err)

	testTracer.Stop()

	spans := testTracer.FinishedSpans()
	require.Equal(t, 1, len(spans))
	finishedSpan := spans[0]
	assert.Equal(t, strconv.Itoa(int(finishedSpan.TraceID())), server.ddTraceID)
	assert.Equal(t, strconv.Itoa(int(finishedSpan.SpanID())), server.ddParentID)

	assert.Empty(t, server.traceparent, "Datadog's mocktracer does not propagate W3C-headers as of writing this test. If they start propagating it, we should remove the separate test below, and update this test to assert the correct W3C-header.")
}

func TestTraceUnaryClientInterceptorW3C(t *testing.T) {
	ctx := context.Background()

	// Ensure valid datadog config, even if we don't have a datadog agent running, to fully instrument the application.
	t.Setenv("DD_ENV", "unittest")
	require.True(t, internal.IsDatadogConfigured())
	t.Setenv("DD_EXPERIMENTAL_TRACING_ENABLED", "true")
	require.True(t, internal.IsExperimentalTracingEnabled())

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

	conn, err := grpc.DialContext(ctx, "bufnet",
		grpc.WithContextDialer(func(ctx context.Context, s string) (net.Conn, error) { return listener.Dial() }),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithUnaryInterceptor(TraceUnaryClientInterceptor()),
	)
	require.NoError(t, err)

	client := testgrpc.NewTestServiceClient(conn)
	span, spanCtx := tracer.StartSpanFromContext(context.Background(), "grpc.request", tracer.ResourceName("/helloworld"))
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
	assert.False(t, strings.Contains(server.tracestate, ","), "w3c tracestate contained multiple list-members, but we only expected 1")
	parts = strings.Split(server.tracestate, "=")
	require.Equal(t, 2, len(parts))
	require.Equal(t, "dd", parts[0], "w3c tracestate list-member did not contain Datadog list-member")

	// Assert Datadog-part of TraceState
	parts = strings.Split(parts[1], ";")
	sort.Strings(parts) // The Datadog-part of the Tracestate can be in any order.
	require.Equal(t, 3, len(parts))
	assert.Equal(t, "s:1", parts[0])
	assert.Equal(t, "t.dm:-1", parts[1])
	assert.True(t, strings.HasPrefix(parts[2], "t.tid:"))
	assert.Equal(t, len("t.tid:65796a3f00000000"), len(parts[2]), "t.tid had invalid length.\nExample: %s\nGot:     %s", "t.tid:65796a3f00000000", parts[2])
}
