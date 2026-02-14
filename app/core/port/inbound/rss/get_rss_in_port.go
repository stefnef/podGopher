package rss

type GetRSSCommand struct {
	ShowSlug         string
	DistributionSlug string
}

type GetRSSResponse struct {
	ShowTitle         string
	DistributionTitle string
	Episodes          []*GetRSSEpisodeResponse
}
type GetRSSEpisodeResponse struct {
	Title string
}

type GetRSSPort interface {
	GetRSS(command *GetRSSCommand) (*GetRSSResponse, error)
}
