package distribution

import (
	"net/http"
	"podGopher/core/port/inbound"
	"podGopher/core/port/inbound/distribution"
	"podGopher/integration/web/handler"

	"github.com/gin-gonic/gin"
)

type CreateDistributionHandler struct {
	route *handler.Route
	port  distribution.CreateDistributionPort
}

type CreateDistributionRequestDto struct {
	Title string `json:"title" binding:"required"`
	Slug  string `json:"slug" binding:"required"`
}

type distributionResponseDto struct {
	Id       string   `json:"id" binding:"required"`
	Title    string   `json:"title" binding:"required"`
	Slug     string   `json:"slug" binding:"required"`
	ShowId   string   `json:"showId" binding:"required"`
	Episodes []string `json:"episodes"`
}

func (h *CreateDistributionHandler) GetRoute() *handler.Route {
	return h.route
}

func NewCreateDistributionHandler(portMap inbound.PortMap) *CreateDistributionHandler {
	return &CreateDistributionHandler{
		route: &handler.Route{
			Method: http.MethodPost,
			Path:   "/show/:showId/distribution",
		},
		port: portMap[inbound.CreateDistribution].(distribution.CreateDistributionPort),
	}
}

func (h *CreateDistributionHandler) Handle(context *gin.Context) {
	showId := context.Param("showId")
	if showId == "" {
		context.AbortWithStatus(http.StatusBadRequest)
		return
	}

	var request *CreateDistributionRequestDto
	if err := context.BindJSON(&request); err != nil {
		context.Abort()
		return
	}

	h.handleCreateDistribution(context, showId, request)
}

func (h *CreateDistributionHandler) handleCreateDistribution(context *gin.Context, showId string, request *CreateDistributionRequestDto) {
	command := &distribution.CreateDistributionCommand{
		ShowId: showId,
		Title:  request.Title,
		Slug:   request.Slug,
	}

	if createdDistribution, err := h.port.CreateDistribution(command); err != nil {
		_ = context.Error(err)
	} else {
		responseDto := distributionResponseDto{
			Id:       createdDistribution.Id,
			Title:    createdDistribution.Title,
			Slug:     createdDistribution.Slug,
			ShowId:   createdDistribution.ShowId,
			Episodes: []string{},
		}
		context.JSON(http.StatusCreated, responseDto)
	}
}
