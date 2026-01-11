package show

import (
	domainError "podGopher/core/domain/error"
	"podGopher/core/domain/model"
	onCreateShow "podGopher/core/port/inbound/show"
	forSaveShow "podGopher/core/port/outbound/show"

	"github.com/google/uuid"
)

type CreateShowService struct {
	saveShowPort forSaveShow.SaveShowPort
}

func NewCreateShowService(repository forSaveShow.SaveShowPort) *CreateShowService {
	return &CreateShowService{
		saveShowPort: repository,
	}
}

func (service *CreateShowService) CreateShow(command *onCreateShow.CreateShowCommand) (*onCreateShow.CreateShowResponse, error) {
	if exists := service.saveShowPort.ExistsByTitleOrSlug(command.Title, command.Slug); exists != false {
		return nil, domainError.NewShowAlreadyExistsError(command.Title)
	}
	id := uuid.NewString()
	show := &model.Show{Id: id, Title: command.Title, Slug: command.Slug}
	err := service.saveShowPort.SaveShow(show)
	if err != nil {
		return nil, err
	}
	return &onCreateShow.CreateShowResponse{Id: show.Id, Title: show.Title, Slug: show.Slug}, nil
}
