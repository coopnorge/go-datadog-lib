package grpc

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/testing/testpb"
)

func TestTraceUnaryServerInterceptor(t *testing.T) {
	grpcUnaryMW := TraceUnaryServerInterceptor()
	grpcUnaryHandler := func(ctx context.Context, req interface{}) (interface{}, error) {
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

func (s *streamServerInterceptorTestSuite) TestPingStream() {
	ctx := context.Background()
	s.Client.PingList(ctx, &testpb.PingListRequest{})
}

func TestTraceStreamServerInterceptor(t *testing.T) {
	s := &streamServerInterceptorTestSuite{
		&testpb.InterceptorTestSuite{
			TestService: &testPingService{
				t: t,
			},
		},
	}
	s.InterceptorTestSuite.ServerOpts = []grpc.ServerOption{
		grpc.StreamInterceptor(TraceStreamServerInterceptor()),
	}
	suite.Run(t, s)
}
