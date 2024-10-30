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
	"sync/atomic"
	"testing"
	"time"

	"github.com/coopnorge/go-datadog-lib/v2/config"
	grpcMiddleware "github.com/coopnorge/go-datadog-lib/v2/middleware/grpc"
	httpMiddleware "github.com/coopnorge/go-datadog-lib/v2/middleware/http"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	testgrpc "google.golang.org/grpc/interop/grpc_testing"
	"google.golang.org/grpc/test/bufconn"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

func TestGRPCServerHttpClientTracing(t *testing.T) {
	// This test aims to ensure that spans are created on incoming gRPC-requests, and then child-spans are created for each outgoing HTTP request, and included in their HTTP Headers.

	// Ensure valid datadog config, even if we don't have a datadog agent running, to fully instrument the application.
	t.Setenv("DD_ENV", "unittest")
	t.Setenv("DD_SERVICE", "unittest")
	t.Setenv("DD_VERSION", "unittest")
	t.Setenv("DD_TRACE_AGENT_URL", "/dev/null")
	t.Setenv("DD_EXPERIMENTAL_TRACING_ENABLED", "true")

	datadogCfg := config.LoadDatadogConfigFromEnvVars()
	require.NoError(t, datadogCfg.Validate())

	// Start Datadog tracer
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
	testgrpc.RegisterTestServiceServer(grpcServer, newStubHealthServer(s.URL))

	// Create in-memory buffer connection and listener
	buffer := 1024 * 1024
	listener := bufconn.Listen(buffer)

	errCh := make(chan error, 1)
	go func() {
		if err := grpcServer.Serve(listener); err != nil && !errors.Is(err, grpc.ErrServerStopped) {
			errCh <- err
			return
		}
	}()

	go func() {
		<-ctx.Done()
		close(errCh)
	}()

	conn, err := grpc.NewClient("dns:///localhost",
		// WithContextDialer connects the client directly to the server over the buffered connection
		grpc.WithContextDialer(func(context.Context, string) (net.Conn, error) { return listener.Dial() }),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	require.NoError(t, err)
	defer conn.Close()

	client := testgrpc.NewTestServiceClient(conn)
	response, err := client.UnaryCall(ctx, &testgrpc.SimpleRequest{})
	require.NoError(t, err)
	assert.Equal(t, "OK", string(response.GetPayload().GetBody()))

	// Assert that we have 1 trace with 3 different spans:
	assert.Equal(t, 3, httpRequestCounter)
	assert.Equal(t, 1, len(traceparentFirst36))
	assert.Equal(t, 3, len(traceparents))
	assert.Equal(t, 1, len(ddTraceIDs))
	assert.Equal(t, 3, len(ddParentIDs)) // note: parent-id is span-id

	// Make second call
	response, err = client.UnaryCall(ctx, &testgrpc.SimpleRequest{})
	require.NoError(t, err)
	assert.Equal(t, "OK", string(response.GetPayload().GetBody()))

	// Assert that we now have 2 traces, and total 6 different spans:
	assert.Equal(t, 6, httpRequestCounter)
	assert.Equal(t, 6, len(traceparents))
	assert.Equal(t, 2, len(traceparentFirst36))
	assert.Equal(t, 6, len(ddParentIDs))
	assert.Equal(t, 2, len(ddTraceIDs))

	cancel() // Signals the gRPC server to stop
	err = <-errCh
	require.NoError(t, err)
}

type testServer struct {
	testgrpc.UnimplementedTestServiceServer
	baseURL string
	client  *http.Client
}

// newStubHealthServer creates a new healthserver that calls external services with tracing.
func newStubHealthServer(baseURL string) *testServer {
	netClient := &http.Client{Timeout: 3 * time.Second}
	netClient = httpMiddleware.AddTracingToClient(
		netClient,
		httpMiddleware.WithResourceNamer(httpMiddleware.FullURLResourceNamer()),
	)
	return &testServer{client: netClient, baseURL: baseURL}
}

// Check implements grpc_health_v1.HealthServer.
func (h *testServer) UnaryCall(ctx context.Context, _ *testgrpc.SimpleRequest) (*testgrpc.SimpleResponse, error) {
	const numServices int = 3
	wg := &sync.WaitGroup{}
	wg.Add(numServices)

	ready := atomic.Bool{}
	ready.Store(true)

	go func() {
		defer wg.Done()
		if err := h.doHTTPRequest(ctx); err != nil {
			ready.Store(false)
		}
	}()

	go func() {
		defer wg.Done()
		if err := h.doHTTPRequest(ctx); err != nil {
			ready.Store(false)
		}
	}()
	go func() {
		defer wg.Done()
		if err := h.doHTTPRequest(ctx); err != nil {
			ready.Store(false)
		}
	}()

	wg.Wait()

	if ready.Load() {
		return &testgrpc.SimpleResponse{Payload: &testgrpc.Payload{Body: []byte("OK")}}, nil
	}
	return &testgrpc.SimpleResponse{Payload: &testgrpc.Payload{Body: []byte("NOT OK")}}, nil
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
