package show

import "podGopher/core/domain/model"

type GetShowCommand struct {
	Id string
}

type GetShowResponse struct {
	Id       string
	Title    string
	Slug     string
	Episodes []string
}

type GetShowPort interface {
	GetShow(command *GetShowCommand) (show *GetShowResponse, err error)
	GetAllShows() (shows *GetAllShowsResponse, err error)
}

type GetAllShowsResponse struct {
	Shows []*model.Show
}
