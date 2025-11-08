package handler

import (
	"net/http"
	"podGopher/core/port/inbound"

	"github.com/gin-gonic/gin"
)

type GetShowHandler struct {
	route *Route
	port  inbound.GetShowPort
}

func NewGetShowHandler(portMap inbound.PortMap) *GetShowHandler {
	return &GetShowHandler{
		route: &Route{
			Method: http.MethodGet,
			Path:   "/show/:showId",
		},
		port: portMap[inbound.GetShow].(inbound.GetShowPort),
	}
}

func (h *GetShowHandler) GetRoute() *Route {
	return h.route
}

func (h *GetShowHandler) Handle(context *gin.Context) {
	foundShow, err := h.port.GetShow(&inbound.GetShowCommand{Id: context.Param("showId")})
	if err != nil {
		_ = context.Error(err)
	} else {
		responseDto := showResponseDto{Id: foundShow.Id, Title: foundShow.Title, Slug: foundShow.Slug}
		context.JSON(http.StatusOK, responseDto)
	}
}
