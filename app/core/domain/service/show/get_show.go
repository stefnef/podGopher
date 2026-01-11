package show

import (
	domainError "podGopher/core/domain/error"
	"podGopher/core/domain/model"
	onGetShow "podGopher/core/port/inbound/show"
	forGetShow "podGopher/core/port/outbound/show"
)

type GetShowService struct {
	repository forGetShow.GetShowPort
}

func NewGetShowService(repository forGetShow.GetShowPort) *GetShowService {
	return &GetShowService{
		repository: repository,
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
