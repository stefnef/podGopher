package distribution

import (
	error2 "podGopher/core/domain/error"
	"podGopher/core/domain/model"
	"podGopher/core/port/inbound"
	"podGopher/core/port/outbound"
)

type GetDistributionService struct {
	getShowOutPort         outbound.GetShowPort
	getDistributionOutPort outbound.GetDistributionPort
}

func NewGetDistributionService(showRepository outbound.GetShowPort, distributionRepository outbound.GetDistributionPort) *GetDistributionService {
	return &GetDistributionService{
		getShowOutPort:         showRepository,
		getDistributionOutPort: distributionRepository,
	}
}

func (service *GetDistributionService) GetDistribution(command *inbound.GetDistributionCommand) (distribution *inbound.GetDistributionResponse, err error) {
	if show, _ := service.getShowOutPort.GetShowOrNil(command.ShowId); show == nil {
		return nil, error2.NewShowNotFoundError(command.ShowId)
	}

	var foundDistribution *model.Distribution
	if foundDistribution, err = service.getDistributionOutPort.GetDistributionOrNil(command.DistributionId); err != nil {
		return nil, err
	}

	if foundDistribution == nil {
		return nil, error2.NewDistributionNotFoundError(command.DistributionId)
	}

	return &inbound.GetDistributionResponse{
		Id:     foundDistribution.Id,
		ShowId: foundDistribution.ShowId,
		Title:  foundDistribution.Title,
		Slug:   foundDistribution.Slug,
	}, nil
}
