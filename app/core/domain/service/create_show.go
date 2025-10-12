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

func (service *CreateShowService) CreateShow(command *inbound.CreateShowCommand) error {
	if exists := service.saveShowPort.ExistsByTitle(command.Title); exists != false {
		return error2.NewShowAlreadyExistsError(command.Title)
	}
	return service.saveShowPort.SaveShow(command.Title)
}
