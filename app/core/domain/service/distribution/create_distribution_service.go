package distribution

import (
	error2 "podGopher/core/domain/error"
	"podGopher/core/domain/model"
	"podGopher/core/port/inbound"
	"podGopher/core/port/outbound"

	"github.com/google/uuid"
)

type CreateDistributionService struct {
	getShowOutPort          outbound.GetShowPort
	saveDistributionOutPort outbound.SaveDistributionPort
}

func NewCreateDistributionService(showRepository outbound.GetShowPort, distributionRepository outbound.SaveDistributionPort) *CreateDistributionService {
	return &CreateDistributionService{
		getShowOutPort:          showRepository,
		saveDistributionOutPort: distributionRepository,
	}
}

func (service CreateDistributionService) CreateDistribution(command *inbound.CreateDistributionCommand) (*inbound.CreateDistributionResponse, error) {
	if exists := service.saveDistributionOutPort.ExistsByTitleOrSlug(command.Title, command.Slug); exists == true {
		return nil, error2.NewDistributionAlreadyExistsError(command.Title)
	}
	if show, _ := service.getShowOutPort.GetShowOrNil(command.ShowId); show == nil {
		return nil, error2.NewShowNotFoundError(command.ShowId)
	}

	id := uuid.NewString()
	distribution := &model.Distribution{Id: id, ShowId: command.ShowId, Title: command.Title, Slug: command.Slug}
	err := service.saveDistributionOutPort.SaveDistribution(distribution)
	if err != nil {
		return nil, err
	}
	return &inbound.CreateDistributionResponse{
		Id:     distribution.Id,
		ShowId: distribution.ShowId,
		Title:  distribution.Title,
		Slug:   distribution.Slug,
	}, nil
}
