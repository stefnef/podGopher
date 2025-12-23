package web

import (
	"errors"
	"net/http"
	error2 "podGopher/core/domain/error"
	"podGopher/core/port/inbound"
	"podGopher/integration/web/handler"
	"podGopher/integration/web/handler/distribution"
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
		episode.NewGetEpisodeHandler(portMap),
		distribution.NewCreateDistributionHandler(portMap),
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
	var showNotFound *error2.ShowNotFoundError
	var episodeAlreadyExists *error2.EpisodeAlreadyExistsError
	var episodeNotFound *error2.EpisodeNotFoundError
	var distributionAlreadyExists *error2.DistributionAlreadyExistsError
	var distributionNotFound *error2.DistributionNotFoundError

	for _, err := range context.Errors {
		switch {
		case errors.As(err.Err, &showAlreadyExists):
			context.AbortWithStatusJSON(http.StatusBadRequest, err.JSON())
		case errors.As(err.Err, &showNotFound):
			context.AbortWithStatusJSON(http.StatusNotFound, err.JSON())
		case errors.As(err.Err, &episodeAlreadyExists):
			context.AbortWithStatusJSON(http.StatusBadRequest, err.JSON())
		case errors.As(err.Err, &episodeNotFound):
			context.AbortWithStatusJSON(http.StatusNotFound, err.JSON())
		case errors.As(err.Err, &distributionAlreadyExists):
			context.AbortWithStatusJSON(http.StatusBadRequest, err.JSON())
		case errors.As(err.Err, &distributionNotFound):
			context.AbortWithStatusJSON(http.StatusNotFound, err.JSON())
		default:
			context.AbortWithStatusJSON(http.StatusInternalServerError, map[string]string{"error": "Internal Server Error"})
		}
	}
	context.Next()
}
