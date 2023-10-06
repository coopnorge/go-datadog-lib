package grpc

import (
	"context"
	"net/http"
	"testing"

	"github.com/coopnorge/go-datadog-lib/v2/internal"
	"github.com/coopnorge/go-datadog-lib/v2/tracing"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/testing/testpb"
)

func TestTraceUnaryServerInterceptorLegacy(t *testing.T) {
	grpcUnaryMW := traceUnaryServerInterceptorLegacy()
	grpcUnaryHandler := func(ctx context.Context, req interface{}) (interface{}, error) {
		meta, exist := internal.GetContextMetadata[tracing.TraceDetails](ctx, internal.TraceContextKey{})
		assert.True(t, exist)
		assert.NotNil(t, meta.DatadogSpan)

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

func TestTraceUnaryServerInterceptorExperimental(t *testing.T) {
	grpcUnaryMW := traceUnaryServerInterceptorExperimental()
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

type testPingServiceLegacy struct {
	*testpb.TestPingService
	t *testing.T
}

func (s *testPingServiceLegacy) PingList(_ *testpb.PingListRequest, stream testpb.TestService_PingListServer) error {
	meta, exist := internal.GetContextMetadata[tracing.TraceDetails](stream.Context(), internal.TraceContextKey{})
	assert.True(s.t, exist)
	assert.NotNil(s.t, meta.DatadogSpan)
	return nil
}

func (s *streamServerInterceptorTestSuite) TestPingStream() {
	ctx := context.Background()
	s.Client.PingList(ctx, &testpb.PingListRequest{})
}

func TestTraceStreamServerInterceptorLegacy(t *testing.T) {
	s := &streamServerInterceptorTestSuite{
		&testpb.InterceptorTestSuite{
			TestService: &testPingServiceLegacy{
				t: t,
			},
		},
	}
	s.InterceptorTestSuite.ServerOpts = []grpc.ServerOption{
		grpc.StreamInterceptor(traceStreamServerInterceptorLegacy()),
	}
	suite.Run(t, s)
}

type testPingServiceExperimental struct {
	*testpb.TestPingService
	t *testing.T
}

func (s *testPingServiceExperimental) PingList(_ *testpb.PingListRequest, stream testpb.TestService_PingListServer) error {
	span, exists := tracer.SpanFromContext(stream.Context())
	assert.True(s.t, exists)
	assert.NotNil(s.t, span)
	return nil
}

func TestTraceStreamServerInterceptorExperimental(t *testing.T) {
	s := &streamServerInterceptorTestSuite{
		&testpb.InterceptorTestSuite{
			TestService: &testPingServiceExperimental{
				t: t,
			},
		},
	}
	s.InterceptorTestSuite.ServerOpts = []grpc.ServerOption{
		grpc.StreamInterceptor(traceStreamServerInterceptorExperimental()),
	}
	suite.Run(t, s)
}
