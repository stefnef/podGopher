package handler

import (
	"podGopher/core/port/inbound"

	"github.com/gin-gonic/gin"
)

type Route struct {
	Method string
	Path   string
}

type Handler interface {
	GetRoute() *Route
	Handle(context *gin.Context)
}

func CreateHandlers(portMap inbound.PortMap) []Handler {
	return []Handler{
		NewCreateShowHandler(portMap),
		NewGetShowHandler(portMap),
		NewCreateEpisodeHandler(portMap),
	}
}
