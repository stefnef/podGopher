package distribution

import (
	"net/http"
	"podGopher/core/port/inbound"
	"podGopher/integration/web/handler"

	"github.com/gin-gonic/gin"
)

type GetDistributionHandler struct {
	route *handler.Route
	port  inbound.GetDistributionPort
}

func (h *GetDistributionHandler) GetRoute() *handler.Route {
	return h.route
}

func NewGetDistributionHandler(portMap inbound.PortMap) *GetDistributionHandler {
	return &GetDistributionHandler{
		route: &handler.Route{
			Method: http.MethodGet,
			Path:   "/show/:showId/distribution/:distributionId",
		},
		port: portMap[inbound.GetDistribution].(inbound.GetDistributionPort),
	}
}

func (h *GetDistributionHandler) Handle(context *gin.Context) {
	showId := context.Param("showId")
	distributionId := context.Param("distributionId")

	if showId == "" || distributionId == "" {
		context.AbortWithStatus(http.StatusBadRequest)
		return
	}

	h.handleGetDistribution(context, showId, distributionId)
}

func (h *GetDistributionHandler) handleGetDistribution(context *gin.Context, showId string, distributionId string) {
	command := &inbound.GetDistributionCommand{
		ShowId:         showId,
		DistributionId: distributionId,
	}

	if foundDistribution, err := h.port.GetDistribution(command); err != nil {
		_ = context.Error(err)
	} else {
		responseDto := distributionResponseDto{
			Id:       foundDistribution.Id,
			Title:    foundDistribution.Title,
			Slug:     foundDistribution.Slug,
			ShowId:   foundDistribution.ShowId,
			Episodes: episodesToDto(foundDistribution),
		}
		context.JSON(http.StatusOK, responseDto)
	}
}

func episodesToDto(foundShow *inbound.GetDistributionResponse) []string {
	var episodesDto []string
	if episodesDto = foundShow.Episodes; episodesDto == nil {
		episodesDto = []string{}
	}
	return episodesDto
}
