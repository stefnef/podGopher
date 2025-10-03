package service

import (
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
		return inbound.NewShowAlreadyExistsError(command.Title)
	}
	return service.saveShowPort.SaveShow(command.Title)
}
