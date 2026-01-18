package distribution

import (
	"errors"
	"net/http"
	domainError "podGopher/core/domain/error"
	"podGopher/core/port/inbound"
	"podGopher/core/port/inbound/distribution"
	"podGopher/integration/web/handler"

	"github.com/gin-gonic/gin"
)

type UpdateDistributionHandler struct {
	route *handler.Route
	port  distribution.UpdateDistributionPort
}

type UpdateDistributionRequestDto struct {
	Title *string `json:"title"`
	Slug  *string `json:"slug"`
}

func (h *UpdateDistributionHandler) GetRoute() *handler.Route {
	return h.route
}

func NewUpdateDistributionHandler(portMap inbound.PortMap) *UpdateDistributionHandler {
	return &UpdateDistributionHandler{
		route: &handler.Route{
			Method: http.MethodPatch,
			Path:   "/show/:showId/distribution/:distributionId",
		},
		port: portMap[inbound.UpdateDistribution].(distribution.UpdateDistributionPort),
	}
}

func (h *UpdateDistributionHandler) Handle(context *gin.Context) {
	showId := context.Param("showId")
	if showId == "" {
		_ = context.Error(domainError.NewShowNotFoundError(""))
		return
	}

	var request *UpdateDistributionRequestDto
	if err := context.BindJSON(&request); err != nil {
		context.Abort()
		return
	}
	if request.Title == nil && request.Slug == nil {
		_ = context.AbortWithError(http.StatusBadRequest, errors.New("invalid request"))
		return
	}

	h.handleUpdateDistribution(context, showId, request)
}

func (h *UpdateDistributionHandler) handleUpdateDistribution(context *gin.Context, showId string, request *UpdateDistributionRequestDto) {
	command := &distribution.UpdateDistributionCommand{
		ShowId: showId,
		Title:  request.Title,
		Slug:   request.Slug,
	}

	if _, err := h.port.UpdateDistribution(command); err != nil {
		_ = context.Error(err)
	} else {
		context.Status(http.StatusNoContent)
		context.Writer.WriteHeaderNow()
	}
}
