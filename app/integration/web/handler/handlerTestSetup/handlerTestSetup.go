package handlerTestSetup

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func GetTestGinContext(t *testing.T) (*gin.Context, *httptest.ResponseRecorder) {
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	ctx, _ := createContextAndEngine(w)

	if ctx == nil {
		t.Fatal()
	}

	return ctx, w
}

func createContextAndEngine(w *httptest.ResponseRecorder) (*gin.Context, *gin.Engine) {
	ctx, engine := gin.CreateTestContext(w)
	ctx.Request = &http.Request{
		Header: make(http.Header),
	}
	return ctx, engine
}
