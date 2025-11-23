package inbound

type CreateEpisodeCommand struct {
	ShowId string
	Title  string
}

type CreateEpisodeResponse struct {
	Id     string
	ShowId string
	Title  string
}

type CreateEpisodePort interface {
	CreateEpisode(command *CreateEpisodeCommand) (episode *CreateEpisodeResponse, err error)
}
