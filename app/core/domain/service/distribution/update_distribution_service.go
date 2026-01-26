package distribution

import (
	"context"
	onUpdateDistribution "podGopher/core/port/inbound/distribution"
	forSaveDistribution "podGopher/core/port/outbound/distribution"
	forGetShow "podGopher/core/port/outbound/show"
)

type UpdateDistributionService struct {
	getShowOutPort          forGetShow.GetShowPort
	saveDistributionOutPort forSaveDistribution.SaveDistributionPort
}

func NewUpdateDistributionService(getShowPort forGetShow.GetShowPort, saveDistributionPort forSaveDistribution.SaveDistributionPort) *UpdateDistributionService {
	return &UpdateDistributionService{
		getShowOutPort:          getShowPort,
		saveDistributionOutPort: saveDistributionPort,
	}
}

func (service UpdateDistributionService) UpdateDistribution(command *onUpdateDistribution.UpdateDistributionCommand) (*onUpdateDistribution.UpdateDistributionResponse, error) {
	context.TODO()
	panic("implement me")
}
