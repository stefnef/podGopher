package show

import (
	"podGopher/core/port/inbound"
	"podGopher/integration/web/handler"

	"github.com/gin-gonic/gin"
)

type GetAllShowsHandler struct {
}

func (h *GetAllShowsHandler) GetRoute() *handler.Route {
	return nil
}

func (h *GetAllShowsHandler) Handle(context *gin.Context) {
	return
}

func NewGetAllShowsHandler(portMap inbound.PortMap) *GetAllShowsHandler {
	return &GetAllShowsHandler{
		/*route: &handler.Route{
			Method: http.MethodGet,
			Path:   "/show/:showId",
		},
		port: portMap[inbound.GetShow].(show.GetShowPort),*/
	}
}
