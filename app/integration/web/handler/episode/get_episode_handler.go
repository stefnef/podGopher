package episode

import (
	"net/http"
	error2 "podGopher/core/domain/error"
	"podGopher/core/port/inbound"
	"podGopher/core/port/inbound/episode"
	"podGopher/integration/web/handler"

	"github.com/gin-gonic/gin"
)

type GetEpisodeHandler struct {
	route *handler.Route
	port  episode.GetEpisodePort
}

func (h GetEpisodeHandler) GetRoute() *handler.Route {
	return h.route
}

func (h GetEpisodeHandler) Handle(context *gin.Context) {
	showId := context.Param("showId")
	episodeId := context.Param("episodeId")
	if showId == "" {
		_ = context.Error(error2.NewShowNotFoundError(""))
		return
	}

	foundEpisode, err := h.port.GetEpisode(&episode.GetEpisodeCommand{
		EpisodeId: episodeId,
		ShowId:    showId,
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
		port: portMap[inbound.GetEpisode].(episode.GetEpisodePort),
	}
}
