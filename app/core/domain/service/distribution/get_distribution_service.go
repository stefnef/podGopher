package distribution

import (
	error2 "podGopher/core/domain/error"
	"podGopher/core/domain/model"
	onGetDistribution "podGopher/core/port/inbound/distribution"
	forGetDistribution "podGopher/core/port/outbound/distribution"
	forGetShow "podGopher/core/port/outbound/show"
)

type GetDistributionService struct {
	getShowOutPort         forGetShow.GetShowPort
	getDistributionOutPort forGetDistribution.GetDistributionPort
}

func NewGetDistributionService(showRepository forGetShow.GetShowPort, distributionRepository forGetDistribution.GetDistributionPort) *GetDistributionService {
	return &GetDistributionService{
		getShowOutPort:         showRepository,
		getDistributionOutPort: distributionRepository,
	}
}

func (service *GetDistributionService) GetDistribution(command *onGetDistribution.GetDistributionCommand) (distribution *onGetDistribution.GetDistributionResponse, err error) {
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

	return &onGetDistribution.GetDistributionResponse{
		Id:     foundDistribution.Id,
		ShowId: foundDistribution.ShowId,
		Title:  foundDistribution.Title,
		Slug:   foundDistribution.Slug,
	}, nil
}
