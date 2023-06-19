package grpc

import (
	"context"
	"net/http"
	"testing"

	"github.com/coopnorge/go-datadog-lib/v2/internal"
	"github.com/coopnorge/go-datadog-lib/v2/tracing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
)

func TestTraceUnaryServerInterceptor(t *testing.T) {
	grpcUnaryMW := TraceUnaryServerInterceptor()
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
