package web

import (
	"errors"
	"net/http"
	error2 "podGopher/core/domain/error"
	"podGopher/core/port/inbound"
	"podGopher/integration/web/handler"
	"podGopher/integration/web/handler/episode"
	"podGopher/integration/web/handler/show"

	"github.com/gin-gonic/gin"
)

func NewRouter(portMap inbound.PortMap) *gin.Engine {
	router := gin.Default()
	setHandlers(portMap, router)

	_ = router.SetTrustedProxies(nil)

	return router
}

func CreateHandlers(portMap inbound.PortMap) []handler.Handler {
	return []handler.Handler{
		show.NewCreateShowHandler(portMap),
		show.NewGetShowHandler(portMap),
		episode.NewCreateEpisodeHandler(portMap),
	}
}

func setHandlers(portMap inbound.PortMap, router *gin.Engine) {
	var handlers = CreateHandlers(portMap)

	for _, handlerImpl := range handlers {
		route := handlerImpl.GetRoute()
		switch route.Method {
		case http.MethodPost:
			router.POST(route.Path, handlerImpl.Handle, handleError)
		case http.MethodGet:
			router.GET(route.Path, handlerImpl.Handle, handleError)
		}
	}
}

func handleError(context *gin.Context) {
	var showAlreadyExists *error2.ShowAlreadyExistsError
	var episodeAlreadyExists *error2.EpisodeAlreadyExistsError
	var showNotFound *error2.ShowNotFoundError

	for _, err := range context.Errors {
		switch {
		case errors.As(err.Err, &showAlreadyExists):
			context.AbortWithStatusJSON(http.StatusBadRequest, err.JSON())
		case errors.As(err.Err, &showNotFound):
			context.AbortWithStatusJSON(http.StatusNotFound, err.JSON())
		case errors.As(err.Err, &episodeAlreadyExists):
			context.AbortWithStatusJSON(http.StatusBadRequest, err.JSON())
		default:
			context.AbortWithStatusJSON(http.StatusInternalServerError, err.JSON())
		}
	}
	context.Next()
}
