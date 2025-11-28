package episode

import (
	"net/http"
	"podGopher/core/port/inbound"
	"podGopher/integration/web/handler"

	"github.com/gin-gonic/gin"
)

type GetEpisodeHandler struct {
	route *handler.Route
	port  inbound.GetEpisodePort
}

func (h GetEpisodeHandler) GetRoute() *handler.Route {
	return h.route
}

func (h GetEpisodeHandler) Handle(context *gin.Context) {
	foundEpisode, err := h.port.GetEpisode(&inbound.GetEpisodeCommand{
		EpisodeId: context.Param("episodeId"),
		ShowId:    context.Param("showId"),
	})
	if err != nil {
		_ = context.Error(err)
	} else {
		responseDto := episodeResponseDto{Id: foundEpisode.Id, ShowId: foundEpisode.ShowId, Title: foundEpisode.Title}
		context.JSON(http.StatusOK, responseDto)
	}
}

func NewGetEpisodeHandler(portMap inbound.PortMap) handler.Handler {
	return GetEpisodeHandler{
		route: &handler.Route{
			Method: http.MethodGet,
			Path:   "/show/:showId/episode/:episodeId",
		},
		port: portMap[inbound.GetEpisode].(inbound.GetEpisodePort),
	}
}
