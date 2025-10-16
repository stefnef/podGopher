package service

import (
	error2 "podGopher/core/domain/error"
	"podGopher/core/port/inbound"
	"podGopher/core/port/outbound"
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
	if exists := service.saveShowPort.ExistsByTitle(command.Title); exists != false {
		return nil, error2.NewShowAlreadyExistsError(command.Title)
	}
	if err := service.saveShowPort.SaveShow(command.Title); err != nil {
		return nil, err
	}
	return &inbound.CreateShowResponse{Title: command.Title}, nil
}
