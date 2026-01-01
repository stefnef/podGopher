package episode

type GetEpisodeCommand struct {
	EpisodeId string
	ShowId    string
}

type GetEpisodeResponse struct {
	Id     string
	ShowId string
	Title  string
}

type GetEpisodePort interface {
	GetEpisode(command *GetEpisodeCommand) (episode *GetEpisodeResponse, err error)
}
