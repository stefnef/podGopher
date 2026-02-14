package rss

import (
	onGetRSS "podGopher/core/port/inbound/rss"
	forGetRSS "podGopher/core/port/outbound/rss"
)

type GetRSSService struct {
	getRSSOutPort forGetRSS.GetRSSPort
}

func NewGetRSSService(RSSRepository forGetRSS.GetRSSPort) *GetRSSService {
	return &GetRSSService{RSSRepository}
}

func (g *GetRSSService) GetRSS(c *onGetRSS.GetRSSCommand) (*onGetRSS.GetRSSResponse, error) {
	rss, err := g.getRSSOutPort.GetRSSOrNil(c.ShowSlug, c.DistributionSlug)
	if err != nil {
		return nil, err
	}

	response := &onGetRSS.GetRSSResponse{
		ShowTitle:         rss.Show.Title,
		DistributionTitle: rss.Distribution.Title,
		Episodes:          []*onGetRSS.GetRSSEpisodeResponse{},
	}

	for _, episode := range rss.Episodes {
		response.Episodes = append(response.Episodes, &onGetRSS.GetRSSEpisodeResponse{Title: episode.Title})
	}

	return response, nil
}
