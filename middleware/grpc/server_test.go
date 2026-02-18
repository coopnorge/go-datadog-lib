package grpc_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/coopnorge/go-datadog-lib/v2/internal/testhelpers"
	datadogMiddleware "github.com/coopnorge/go-datadog-lib/v2/middleware/grpc"

	"github.com/DataDog/dd-trace-go/v2/ddtrace/mocktracer"
	"github.com/DataDog/dd-trace-go/v2/ddtrace/tracer"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/testing/testpb"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc"
)

func TestTraceUnaryServerInterceptor(t *testing.T) {
	// Start Datadog tracer, so that we don't create NoopSpans.
	testTracer := mocktracer.Start()
	t.Cleanup(testTracer.Stop)

	testhelpers.ConfigureDatadog(t)

	grpcUnaryMW := datadogMiddleware.TraceUnaryServerInterceptor()
	grpcUnaryHandler := func(ctx context.Context, _ interface{}) (interface{}, error) {
		span, exists := tracer.SpanFromContext(ctx)
		assert.True(t, exists)
		assert.NotNil(t, span)

		return nil, nil
	}

	tCtx := context.Background()
	tReq, err := http.NewRequest(http.MethodGet, "unit.test", nil)
	require.NoError(t, err)
	resp, err := grpcUnaryMW(
		tCtx,
		tReq,
		&grpc.UnaryServerInfo{FullMethod: "test"},
		grpcUnaryHandler,
	)

	require.NoError(t, err)
	assert.Nil(t, resp)
}

type streamServerInterceptorTestSuite struct {
	*testpb.InterceptorTestSuite
}

func (s *streamServerInterceptorTestSuite) TestPingStream() {
	ctx := context.Background()
	_, err := s.Client.PingList(ctx, &testpb.PingListRequest{})
	require.NoError(s.T(), err)
}

type testPingService struct {
	*testpb.TestPingService
	t *testing.T
}

func (s *testPingService) PingList(_ *testpb.PingListRequest, stream testpb.TestService_PingListServer) error {
	span, exists := tracer.SpanFromContext(stream.Context())
	assert.True(s.t, exists)
	assert.NotNil(s.t, span)
	return nil
}

func TestTraceStreamServerInterceptor(t *testing.T) {
	// Start Datadog tracer, so that we don't create NoopSpans.
	testTracer := mocktracer.Start()
	t.Cleanup(testTracer.Stop)

	s := &streamServerInterceptorTestSuite{
		InterceptorTestSuite: &testpb.InterceptorTestSuite{
			TestService: &testPingService{&testpb.TestPingService{}, t},
			ServerOpts: []grpc.ServerOption{
				grpc.StreamInterceptor(datadogMiddleware.TraceStreamServerInterceptor()),
			},
		},
	}
	suite.Run(t, s)
}
