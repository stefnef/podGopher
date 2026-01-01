package show

import (
	"net/http"
	"podGopher/core/port/inbound"
	"podGopher/core/port/inbound/show"
	"podGopher/integration/web/handler"

	"github.com/gin-gonic/gin"
)

type GetShowHandler struct {
	route *handler.Route
	port  show.GetShowPort
}

func NewGetShowHandler(portMap inbound.PortMap) *GetShowHandler {
	return &GetShowHandler{
		route: &handler.Route{
			Method: http.MethodGet,
			Path:   "/show/:showId",
		},
		port: portMap[inbound.GetShow].(show.GetShowPort),
	}
}

func (h *GetShowHandler) GetRoute() *handler.Route {
	return h.route
}

func (h *GetShowHandler) Handle(context *gin.Context) {
	foundShow, err := h.port.GetShow(&show.GetShowCommand{Id: context.Param("showId")})
	if err != nil {
		_ = context.Error(err)
	} else {
		responseDto := showResponseDto{Id: foundShow.Id, Title: foundShow.Title, Slug: foundShow.Slug, Episodes: episodesToDto(foundShow)}
		context.JSON(http.StatusOK, responseDto)
	}
}

func episodesToDto(foundShow *show.GetShowResponse) []string {
	var episodesDto []string
	if episodesDto = foundShow.Episodes; episodesDto == nil {
		episodesDto = []string{}
	}
	return episodesDto
}
