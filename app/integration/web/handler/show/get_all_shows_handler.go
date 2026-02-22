package show

import (
	"net/http"
	"podGopher/core/port/inbound"
	"podGopher/core/port/inbound/show"
	"podGopher/integration/web/handler"

	"github.com/gin-gonic/gin"
)

type GetAllShowsHandler struct {
	route *handler.Route
	port  show.GetShowPort
}

func (h *GetAllShowsHandler) GetRoute() *handler.Route {
	return h.route
}

func (h *GetAllShowsHandler) Handle(context *gin.Context) {
	_, err := h.port.GetAllShows()
	if err != nil {
		_ = context.Error(err)
	}
	return
}

func NewGetAllShowsHandler(portMap inbound.PortMap) *GetAllShowsHandler {
	return &GetAllShowsHandler{
		route: &handler.Route{
			Method: http.MethodGet,
			Path:   "/show",
		},
		port: portMap[inbound.GetShow].(show.GetShowPort),
	}
}
