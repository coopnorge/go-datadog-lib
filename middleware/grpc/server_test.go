package grpc

import (
	"context"
	"net/http"
	"testing"

	"github.com/coopnorge/go-datadog-lib/v2/internal/testhelpers"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/testing/testpb"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/mocktracer"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

func TestTraceUnaryServerInterceptor(t *testing.T) {
	// Start Datadog tracer, so that we don't create NoopSpans.
	testTracer := mocktracer.Start()
	t.Cleanup(testTracer.Stop)

	testhelpers.ConfigureDatadog(t)

	grpcUnaryMW := TraceUnaryServerInterceptor()
	grpcUnaryHandler := func(ctx context.Context, _ interface{}) (interface{}, error) {
		span, exists := tracer.SpanFromContext(ctx)
		assert.True(t, exists)
		assert.NotNil(t, span)

		return nil, nil
	}

	tCtx := context.Background()
	tReq, _ := http.NewRequest(http.MethodGet, "unit.test", nil)
	resp, err := grpcUnaryMW(
		tCtx,
		tReq,
		&grpc.UnaryServerInfo{FullMethod: "test"},
		grpcUnaryHandler,
	)

	assert.Nil(t, err)
	assert.Nil(t, resp)
}

type streamServerInterceptorTestSuite struct {
	*testpb.InterceptorTestSuite
}

func (s *streamServerInterceptorTestSuite) TestPingStream() {
	ctx := context.Background()
	s.Client.PingList(ctx, &testpb.PingListRequest{})
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
				grpc.StreamInterceptor(TraceStreamServerInterceptor()),
			},
		},
	}
	suite.Run(t, s)
}
