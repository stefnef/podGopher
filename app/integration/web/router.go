package web

import (
	"errors"
	"net/http"
	error2 "podGopher/core/domain/error"
	"podGopher/core/port/inbound"
	"podGopher/integration/web/handler"

	"github.com/gin-gonic/gin"
)

func NewRouter(portMap inbound.PortMap) *gin.Engine {
	router := gin.Default()
	var handlers = handler.CreateHandlers(portMap)

	for _, handlerImpl := range handlers {
		route := handlerImpl.GetRoute()
		switch route.Method {
		case http.MethodPost:
			router.POST(route.Path, handlerImpl.Handle, handleError)
		}
	}
	return router
}

func handleError(context *gin.Context) {
	var alreadyExists *error2.ShowAlreadyExistsError

	for _, err := range context.Errors {
		switch {
		case errors.As(err.Err, &alreadyExists):
			context.AbortWithStatusJSON(http.StatusBadRequest, err.JSON())
		default:
			context.AbortWithStatusJSON(http.StatusInternalServerError, err.JSON())
		}
	}
	context.Next()
}
