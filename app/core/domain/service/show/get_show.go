package show

import (
	domainError "podGopher/core/domain/error"
	"podGopher/core/domain/model"
	onGetShow "podGopher/core/port/inbound/show"
	forGetShow "podGopher/core/port/outbound/show"
)

type GetShowService struct {
	repository         forGetShow.GetShowPort
	repositoryAllShows forGetShow.GetAllShowsPort
}

func NewGetShowService(repository forGetShow.GetShowPort, repositoryAllShows forGetShow.GetAllShowsPort) *GetShowService {
	return &GetShowService{
		repository:         repository,
		repositoryAllShows: repositoryAllShows,
	}
}

func (s *GetShowService) GetShow(command *onGetShow.GetShowCommand) (showResponse *onGetShow.GetShowResponse, err error) {
	var show *model.Show
	if show, err = s.repository.GetShowOrNil(command.Id); err != nil {
		return nil, err
	}

	if show == nil {
		return nil, domainError.NewShowNotFoundError(command.Id)
	}
	return &onGetShow.GetShowResponse{
		Id:       show.Id,
		Title:    show.Title,
		Slug:     show.Slug,
		Episodes: show.Episodes,
	}, nil
}

func (s *GetShowService) GetAllShows() (shows *onGetShow.GetAllShowsResponse, err error) {
	allShows, err := s.repositoryAllShows.GetAllShows()

	return &onGetShow.GetAllShowsResponse{
		Shows: allShows,
	}, err
}
