package show

import (
	error2 "podGopher/core/domain/error"
	"podGopher/core/domain/model"
	"podGopher/core/port/inbound"
	"podGopher/core/port/outbound"
)

type GetShowService struct {
	repository outbound.GetShowPort
}

func NewGetShowService(repository outbound.GetShowPort) *GetShowService {
	return &GetShowService{
		repository: repository,
	}
}

func (s *GetShowService) GetShow(command *inbound.GetShowCommand) (showResponse *inbound.GetShowResponse, err error) {
	var show *model.Show
	if show, err = s.repository.GetShowOrNil(command.Id); err != nil {
		return nil, err
	}

	if show == nil {
		return nil, error2.NewShowNotFoundError(command.Id)
	}
	return &inbound.GetShowResponse{
		Id:       show.Id,
		Title:    show.Title,
		Slug:     show.Slug,
		Episodes: show.Episodes,
	}, nil
}
