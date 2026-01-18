package distribution

import (
	"context"
	onUpdateDistribution "podGopher/core/port/inbound/distribution"
)

type UpdateDistributionService struct {
}

func NewUpdateDistributionService() *UpdateDistributionService {
	return &UpdateDistributionService{}
}

func (service UpdateDistributionService) UpdateDistribution(command *onUpdateDistribution.UpdateDistributionCommand) (*onUpdateDistribution.UpdateDistributionResponse, error) {
	context.TODO()
	panic("implement me")
}
