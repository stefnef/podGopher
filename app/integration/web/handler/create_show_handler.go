package handler

import (
	"net/http"
	"podGopher/core/port/inbound"

	"github.com/gin-gonic/gin"
)

type CreateShowHandler struct {
	route *Route
	port  inbound.CreateShowPort
}

type CreateShowRequestDto struct {
	Title string `json:"title" binding:"required"`
	Slug  string `json:"slug" binding:"required"`
}

type showResponseDto struct {
	Id    string `json:"id" binding:"required"`
	Title string `json:"title" binding:"required"`
	Slug  string `json:"slug" binding:"required"`
}

func (h *CreateShowHandler) GetRoute() *Route {
	return h.route
}

func NewCreateShowHandler(portMap inbound.PortMap) *CreateShowHandler {
	return &CreateShowHandler{
		route: &Route{
			Method: http.MethodPost,
			Path:   "/show",
		},
		port: portMap[inbound.CreateShow].(inbound.CreateShowPort),
	}
}

func (h *CreateShowHandler) Handle(context *gin.Context) {
	var request *CreateShowRequestDto
	if err := context.BindJSON(&request); err != nil {
		context.Abort()
		return
	}

	h.handleCreateShow(context, request)
}

func (h *CreateShowHandler) handleCreateShow(context *gin.Context, request *CreateShowRequestDto) {
	if createdShow, err := h.port.CreateShow(&inbound.CreateShowCommand{Title: request.Title, Slug: request.Slug}); err != nil {
		_ = context.Error(err)
	} else {
		responseDto := showResponseDto{Id: createdShow.Id, Title: createdShow.Title, Slug: createdShow.Slug}
		context.JSON(http.StatusAccepted, responseDto)
	}
}
