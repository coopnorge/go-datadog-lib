package echo

import (
	"net/http"
	"testing"

	"github.com/coopnorge/go-datadog-lib/internal"
	"github.com/coopnorge/go-datadog-lib/internal/generated/mocks/labstack/echo/v4"
	"github.com/coopnorge/go-datadog-lib/tracing"

	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestTraceServerMiddleware(t *testing.T) {
	echoMiddlewareHandler := TraceServerMiddleware()
	echoRequestHandler := func(reqCtx echo.Context) (err error) {
		assert.NotNil(t, reqCtx.Request())
		// Since there is mock you cannot fetch TraceDetails to verify it
		return nil
	}

	tReq, _ := http.NewRequest(http.MethodGet, "unit.test", nil)

	ctrl := gomock.NewController(t)
	mockEchoContext := mock_echo.NewMockContext(ctrl)
	ctrl.Finish()

	mockEchoContext.EXPECT().Request().Return(tReq).MaxTimes(5)
	mockEchoContext.EXPECT().SetRequest(gomock.Any()).MaxTimes(1)

	echoMiddlewareFun := echoMiddlewareHandler(echoRequestHandler)
	echoHandlerFunc := echoMiddlewareFun(mockEchoContext)

	assert.Nil(t, echoHandlerFunc)
}

func TestTraceServerMiddlewareEchoMissingRequest(t *testing.T) {
	echoMiddlewareHandler := TraceServerMiddleware()
	echoRequestHandler := func(reqCtx echo.Context) (err error) {
		assert.NotNil(t, reqCtx.Request())

		meta, exist := internal.GetContextMetadata[tracing.TraceDetails](reqCtx.Request().Context(), internal.TraceContextKey{})
		assert.True(t, exist)
		assert.NotNil(t, meta.DatadogSpan)

		return nil
	}

	ctrl := gomock.NewController(t)
	mockEchoContext := mock_echo.NewMockContext(ctrl)
	ctrl.Finish()

	mockEchoContext.EXPECT().Request().Return(nil).MaxTimes(1)

	echoMiddlewareFun := echoMiddlewareHandler(echoRequestHandler)
	echoHandlerFunc := echoMiddlewareFun(mockEchoContext)

	assert.NotNil(t, echoHandlerFunc)
}
