package service

import (
	error2 "podGopher/core/domain/error"
	"podGopher/core/domain/model"
	"podGopher/core/port/inbound"
	"podGopher/core/port/outbound"

	"github.com/google/uuid"
)

type CreateShowService struct {
	saveShowPort outbound.SaveShowPort
}

func NewCreateShowService(repository outbound.SaveShowPort) *CreateShowService {
	return &CreateShowService{
		saveShowPort: repository,
	}
}

func (service *CreateShowService) CreateShow(command *inbound.CreateShowCommand) (*inbound.CreateShowResponse, error) {
	if exists := service.saveShowPort.ExistsByTitleOrSlug(command.Title, command.Slug); exists != false {
		return nil, error2.NewShowAlreadyExistsError(command.Title)
	}
	id := uuid.NewString()
	show := &model.Show{Id: id, Title: command.Title, Slug: command.Slug}
	err := service.saveShowPort.SaveShow(show)
	if err != nil {
		return nil, err
	}
	return &inbound.CreateShowResponse{Id: show.Id, Title: show.Title, Slug: show.Slug}, nil
}
