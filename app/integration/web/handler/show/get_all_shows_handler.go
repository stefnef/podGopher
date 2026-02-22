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

type allShowsResponseDto struct {
	Shows []allShowsItemResponseDto `json:"shows" binding:"required"`
}

type allShowsItemResponseDto struct {
	Id    string `json:"id" binding:"required"`
	Title string `json:"title" binding:"required"`
}

func (h *GetAllShowsHandler) GetRoute() *handler.Route {
	return h.route
}

func (h *GetAllShowsHandler) Handle(context *gin.Context) {
	foundShows, err := h.port.GetAllShows()
	if err != nil {
		_ = context.Error(err)
	} else {
		responseDto := mapAllFoundShows(foundShows)
		context.JSON(http.StatusOK, responseDto)
	}
}

func mapAllFoundShows(foundShows *show.GetAllShowsResponse) allShowsResponseDto {
	responseDto := allShowsResponseDto{[]allShowsItemResponseDto{}}
	for _, singleShow := range foundShows.Shows {
		responseDto.Shows = append(responseDto.Shows, allShowsItemResponseDto{
			Id:    singleShow.Id,
			Title: singleShow.Title,
		})
	}
	return responseDto
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
