package echo_test

import (
	"net/http"
	"testing"

	mock_echo "github.com/coopnorge/go-datadog-lib/v2/internal/generated/mocks/labstack/echo/v4"
	"github.com/coopnorge/go-datadog-lib/v2/internal/testhelpers"

	coopEchoDatadog "github.com/coopnorge/go-datadog-lib/v2/middleware/echo"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	gomock "go.uber.org/mock/gomock"
)

func TestTraceServerMiddleware(t *testing.T) {
	testhelpers.ConfigureDatadog(t)

	echoMiddlewareHandler := coopEchoDatadog.TraceServerMiddleware()
	echoRequestHandler := func(reqCtx echo.Context) (err error) {
		assert.NotNil(t, reqCtx.Request())
		// Since there is mock you cannot fetch TraceDetails to verify it
		return nil
	}

	tReq, err := http.NewRequest(http.MethodGet, "unit.test", nil)
	require.NoError(t, err)

	ctrl := gomock.NewController(t)
	mockEchoContext := mock_echo.NewMockContext(ctrl)

	mockEchoContext.EXPECT().Request().Return(tReq).MaxTimes(5)
	mockEchoContext.EXPECT().SetRequest(gomock.Any()).MaxTimes(1)
	mockEchoContext.EXPECT().Path().Return("")
	mockEchoContext.EXPECT().Response().Return(&echo.Response{})

	echoMiddlewareFun := echoMiddlewareHandler(echoRequestHandler)
	echoHandlerFunc := echoMiddlewareFun(mockEchoContext)

	assert.Nil(t, echoHandlerFunc)
}
