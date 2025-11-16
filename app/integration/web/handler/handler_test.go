package handler

import (
	"net/http"
	"net/http/httptest"
	"podGopher/core/domain/service"
	"podGopher/core/port/inbound"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func Test_should_create_handlers(t *testing.T) {
	portMap := inbound.PortMap{
		inbound.CreateShow:    service.NewCreateShowService(nil),
		inbound.GetShow:       service.NewGetShowService(nil),
		inbound.CreateEpisode: service.NewCreateEpisodeService(nil, nil),
	}

	var handlers = CreateHandlers(portMap)

	assert.NotEmpty(t, handlers)
	assert.Len(t, handlers, 3)
}

func GetTestGinContext() (*gin.Context, *httptest.ResponseRecorder) {
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	ctx, _ := createContextAndEngine(w)

	return ctx, w
}

func createContextAndEngine(w *httptest.ResponseRecorder) (*gin.Context, *gin.Engine) {
	ctx, engine := gin.CreateTestContext(w)
	ctx.Request = &http.Request{
		Header: make(http.Header),
	}
	return ctx, engine
}
