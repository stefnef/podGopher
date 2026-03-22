package rss

import (
	"net/http"
	domainError "podGopher/core/domain/error"
	"podGopher/core/port/inbound"
	"podGopher/core/port/inbound/rss"
	"podGopher/integration/web/handler"

	"github.com/gin-gonic/gin"
)

type rssResponseDto struct {
	ShowTitle         string `json:"showTitle" binding:"required"`
	DistributionTitle string `json:"DistributionTitle" binding:"required"`
	Episodes          []rssEpisodeResponseDto
}

type rssEpisodeResponseDto struct {
	Title string `json:"title" binding:"required"`
}

type GetRSSHandler struct {
	route *handler.Route
	port  rss.GetRSSPort
}

func (h GetRSSHandler) Authorize(*gin.Context) {
	//TODO implement me
}

func NewGetRSSHandler(portMap inbound.PortMap) handler.Handler {
	return &GetRSSHandler{
		route: &handler.Route{
			Method: http.MethodGet,
			Path:   "/rss/:showSlug/:distributionSlug",
		},
		port: portMap[inbound.GetRSS].(rss.GetRSSPort),
	}
}

func (h GetRSSHandler) GetRoute() *handler.Route {
	return h.route
}

func (h GetRSSHandler) Handle(context *gin.Context) {
	showSlug := context.Param("showSlug")
	distributionSlug := context.Param("distributionSlug")
	if showSlug == "" || distributionSlug == "" {
		_ = context.Error(domainError.NewRSSFeedNotFoundError(showSlug + distributionSlug))
		return
	}

	foundRSS, err := h.port.GetRSS(&rss.GetRSSCommand{
		ShowSlug:         showSlug,
		DistributionSlug: distributionSlug,
	})
	if err != nil {
		_ = context.Error(err)
	} else {
		responseDto := rssResponseDto{
			ShowTitle:         foundRSS.ShowTitle,
			DistributionTitle: foundRSS.DistributionTitle,
			Episodes:          []rssEpisodeResponseDto{},
		}
		for _, episode := range foundRSS.Episodes {
			responseDto.Episodes = append(responseDto.Episodes, rssEpisodeResponseDto{Title: episode.Title})
		}
		context.XML(http.StatusOK, responseDto)
	}
}
