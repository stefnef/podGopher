package handler

import (
	"net/http"
	"podGopher/core/port/inbound"

	"github.com/gin-gonic/gin"
)

type CreateEpisodeHandler struct {
	route *Route
	port  inbound.CreateEpisodePort
}

type CreateEpisodeRequestDto struct {
	Title string `json:"title" binding:"required"`
}

type episodeResponseDto struct {
	Id     string `json:"id" binding:"required"`
	ShowId string `json:"showId" binding:"required"`
	Title  string `json:"title" binding:"required"`
}

func (h *CreateEpisodeHandler) GetRoute() *Route {
	return h.route
}

func NewCreateEpisodeHandler(portMap inbound.PortMap) *CreateEpisodeHandler {
	return &CreateEpisodeHandler{
		route: &Route{
			Method: http.MethodPost,
			Path:   "/show/:showId/episode",
		},
		port: portMap[inbound.CreateEpisode].(inbound.CreateEpisodePort),
	}
}

func (h *CreateEpisodeHandler) Handle(context *gin.Context) {
	var request *CreateEpisodeRequestDto
	if err := context.BindJSON(&request); err != nil {
		context.Abort()
		return
	}

	h.handleCreateEpisode(context, request)
}

func (h *CreateEpisodeHandler) handleCreateEpisode(context *gin.Context, request *CreateEpisodeRequestDto) {
	createEpisodeCommand := &inbound.CreateEpisodeCommand{ShowId: context.Param("showId"), Title: request.Title}
	if createdEpisode, err := h.port.CreateEpisode(createEpisodeCommand); err != nil {
		_ = context.Error(err)
	} else {
		responseDto := episodeResponseDto{Id: createdEpisode.Id, ShowId: createdEpisode.ShowId, Title: createdEpisode.Title}
		context.JSON(http.StatusCreated, responseDto)
	}
}
