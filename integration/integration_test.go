package integration

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	"github.com/coopnorge/go-datadog-lib/v2/internal/testhelpers"
	grpcMiddleware "github.com/coopnorge/go-datadog-lib/v2/middleware/grpc"
	httpMiddleware "github.com/coopnorge/go-datadog-lib/v2/middleware/http"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	testgrpc "google.golang.org/grpc/interop/grpc_testing"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/test/bufconn"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

const bufSize = 1024 * 1024

func TestGRPCServerHttpClientTracing(t *testing.T) {
	// This test aims to ensure that spans are created on incoming gRPC-requests, and then child-spans are created for each outgoing HTTP request, and included in their HTTP Headers.

	testhelpers.ConfigureDatadog(t)

	// Start Datadog tracer, so that we don't create NoopSpans.
	// Start real tracer (not mocktracer), to propagate Traceparent.
	tracer.Start(tracer.WithService("unittest"))
	t.Cleanup(tracer.Flush)
	t.Cleanup(tracer.Stop)

	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)

	mu := sync.Mutex{}
	traceparentFirst36 := make(map[string]struct{})
	traceparents := make(map[string]struct{})
	ddParentIDs := make(map[string]struct{})
	ddTraceIDs := make(map[string]struct{})
	httpRequestCounter := 0
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		defer mu.Unlock()
		httpRequestCounter++
		traceparentFirst36[r.Header.Get("Traceparent")[0:36]] = struct{}{} // Grab the first 36 runes, which contain the trace-id
		traceparents[r.Header.Get("Traceparent")] = struct{}{}
		ddTraceIDs[r.Header.Get("X-Datadog-Trace-Id")] = struct{}{}
		ddParentIDs[r.Header.Get("X-Datadog-Parent-Id")] = struct{}{}
		_, _ = w.Write([]byte("OK"))
	}))
	t.Cleanup(s.Close)

	serverOpts := []grpc.ServerOption{
		grpc.UnaryInterceptor(grpcMiddleware.UnaryServerInterceptor()),
	}
	grpcServer := grpc.NewServer(serverOpts...)
	testgrpc.RegisterTestServiceServer(grpcServer, newTestServer(s.URL, 3))

	listener := bufconn.Listen(bufSize)

	errCh := make(chan error, 1)
	go func() {
		errCh <- grpcServer.Serve(listener)
	}()

	conn, err := grpc.NewClient("dns:///localhost",
		grpc.WithContextDialer(func(context.Context, string) (net.Conn, error) { return listener.Dial() }),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithUnaryInterceptor(grpcMiddleware.UnaryClientInterceptor()),
	)
	require.NoError(t, err)
	defer conn.Close()

	client := testgrpc.NewTestServiceClient(conn)
	_, err = client.EmptyCall(ctx, &testgrpc.Empty{})
	require.NoError(t, err)

	// Assert that we have 1 trace with 3 different spans:
	assert.Equal(t, 3, httpRequestCounter)
	assert.Equal(t, 1, len(traceparentFirst36))
	assert.Equal(t, 3, len(traceparents))
	assert.Equal(t, 1, len(ddTraceIDs))
	assert.Equal(t, 3, len(ddParentIDs)) // note: X-Datadog-Parent-Id is span-id

	// Make second call
	_, err = client.EmptyCall(ctx, &testgrpc.Empty{})
	require.NoError(t, err)

	// Assert that we now have 2 traces, and total 6 different spans:
	assert.Equal(t, 6, httpRequestCounter)
	assert.Equal(t, 2, len(traceparentFirst36))
	assert.Equal(t, 6, len(traceparents))
	assert.Equal(t, 2, len(ddTraceIDs))
	assert.Equal(t, 6, len(ddParentIDs)) // note: X-Datadog-Parent-Id is span-id

	grpcServer.Stop()
	err = <-errCh
	require.NoError(t, err)
}

type testServer struct {
	testgrpc.UnimplementedTestServiceServer

	baseURL          string
	client           *http.Client
	numExternalCalls int
}

// newTestServer creates a new server that implements testgrpc.TestServiceServer that calls external services via HTTP with tracing.
func newTestServer(baseURL string, numExternalCalls int) *testServer {
	netClient := &http.Client{Timeout: 3 * time.Second}
	netClient = httpMiddleware.AddTracingToClient(
		netClient,
		httpMiddleware.WithResourceNamer(httpMiddleware.FullURLResourceNamer()),
	)
	return &testServer{client: netClient, baseURL: baseURL, numExternalCalls: numExternalCalls}
}

// Check implements testgrpc.TestServiceServer.
func (h *testServer) EmptyCall(ctx context.Context, _ *testgrpc.Empty) (*testgrpc.Empty, error) {
	g, ctx := errgroup.WithContext(ctx)
	for i := 0; i < h.numExternalCalls; i++ {
		g.Go(func() error { return h.doHTTPRequest(ctx) })
	}

	err := g.Wait()
	if err != nil {
		return nil, status.Errorf(codes.Unavailable, "Service is not ready")
	}

	return &testgrpc.Empty{}, nil
}

func (h *testServer) doHTTPRequest(ctx context.Context) (err error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, h.baseURL, nil)
	if err != nil {
		return fmt.Errorf("error creating http.Request: %w", err)
	}

	resp, err := h.client.Do(req)
	if err != nil {
		return fmt.Errorf("unable to send request: %w", err)
	}
	defer func() {
		if closeErr := resp.Body.Close(); closeErr != nil {
			// Mutate the named return, to capture any errors during closing of response-body
			err = errors.Join(err, closeErr)
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code returned: %d", resp.StatusCode)
	}

	// drain the response body to complete the request.
	_, err = io.Copy(io.Discard, resp.Body)
	if err != nil {
		return fmt.Errorf("error reading(discarding) body: %w", err)
	}

	return nil
}
