package echo

import (
	"net/http"
	"testing"

	mock_echo "github.com/coopnorge/go-datadog-lib/v2/internal/generated/mocks/labstack/echo/v4"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"

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
	mockEchoContext.EXPECT().Path().Return("")
	mockEchoContext.EXPECT().Response().Return(&echo.Response{})

	echoMiddlewareFun := echoMiddlewareHandler(echoRequestHandler)
	echoHandlerFunc := echoMiddlewareFun(mockEchoContext)

	assert.Nil(t, echoHandlerFunc)
}

func TestTraceServerMiddlewareEchoMissingRequest(t *testing.T) {
	echoMiddlewareHandler := TraceServerMiddleware()
	echoRequestHandler := func(reqCtx echo.Context) (err error) {
		assert.NotNil(t, reqCtx.Request())

		span, exists := tracer.SpanFromContext(reqCtx.Request().Context())
		assert.True(t, exists)
		assert.NotNil(t, span)

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
