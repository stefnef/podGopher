package web

import (
	"errors"
	"net/http"
	domainError "podGopher/core/domain/error"
	"podGopher/core/port/inbound"
	"podGopher/integration/web/auth"
	"podGopher/integration/web/handler"
	"podGopher/integration/web/handler/distribution"
	"podGopher/integration/web/handler/episode"
	"podGopher/integration/web/handler/rss"
	"podGopher/integration/web/handler/show"
	"podGopher/integration/web/handler/user"

	"github.com/gin-gonic/gin"
)

func NewRouter(portMap inbound.PortMap, adminAuth auth.AdminAuth) *gin.Engine {
	router := gin.Default()
	setHandlers(portMap, router, adminAuth)

	_ = router.SetTrustedProxies(nil)

	return router
}

func CreateHandlers(portMap inbound.PortMap, adminAuth auth.AdminAuth) []handler.Handler {
	return []handler.Handler{
		show.NewCreateShowHandler(portMap),
		show.NewGetShowHandler(portMap),
		show.NewGetAllShowsHandler(portMap),
		episode.NewCreateEpisodeHandler(portMap),
		episode.NewGetEpisodeHandler(portMap),
		distribution.NewCreateDistributionHandler(portMap),
		distribution.NewGetDistributionHandler(portMap),
		distribution.NewUpdateDistributionHandler(portMap),
		rss.NewGetRSSHandler(portMap),
		user.NewCreateUserHandler(portMap, adminAuth),
	}
}

func setHandlers(portMap inbound.PortMap, router *gin.Engine, adminAuth auth.AdminAuth) {
	var handlers = CreateHandlers(portMap, adminAuth)

	for _, handlerImpl := range handlers {
		route := handlerImpl.GetRoute()
		switch route.Method {
		case http.MethodPost:
			router.POST(route.Path, handlerImpl.Authorize, handlerImpl.Handle, handleError)
		case http.MethodGet:
			router.GET(route.Path, handlerImpl.Handle, handleError)
		case http.MethodPatch:
			router.PATCH(route.Path, handlerImpl.Handle, handleError)
		}
	}
}

func handleError(context *gin.Context) {
	var modelError *domainError.ModelError

	for _, err := range context.Errors {
		switch {
		case errors.As(err.Err, &modelError):
			switch err.Err.(*domainError.ModelError).Category {
			case domainError.AlreadyExists:
				context.AbortWithStatusJSON(http.StatusBadRequest, err.JSON())
			case domainError.NotFound:
				context.AbortWithStatusJSON(http.StatusNotFound, err.JSON())
			case domainError.DataConflict:
				context.AbortWithStatusJSON(http.StatusBadRequest, err.JSON())
			case domainError.Authorization:
				context.AbortWithStatusJSON(http.StatusForbidden, err.JSON())
			default:
				context.AbortWithStatusJSON(http.StatusInternalServerError, err.JSON())
			}
		default:
			context.AbortWithStatusJSON(http.StatusInternalServerError, map[string]string{"error": "Internal Server Error"})
		}
	}
	context.Next()
}
