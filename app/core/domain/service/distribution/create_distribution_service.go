package distribution

import (
	domainError "podGopher/core/domain/error"
	"podGopher/core/domain/model"
	onCreateDistribution "podGopher/core/port/inbound/distribution"
	forSaveDistribution "podGopher/core/port/outbound/distribution"
	forGetShow "podGopher/core/port/outbound/show"

	"github.com/google/uuid"
)

type CreateDistributionService struct {
	getShowOutPort          forGetShow.GetShowPort
	saveDistributionOutPort forSaveDistribution.SaveDistributionPort
}

func NewCreateDistributionService(showRepository forGetShow.GetShowPort, distributionRepository forSaveDistribution.SaveDistributionPort) *CreateDistributionService {
	return &CreateDistributionService{
		getShowOutPort:          showRepository,
		saveDistributionOutPort: distributionRepository,
	}
}

func (service CreateDistributionService) CreateDistribution(command *onCreateDistribution.CreateDistributionCommand) (*onCreateDistribution.CreateDistributionResponse, error) {
	if exists := service.saveDistributionOutPort.ExistsByTitleOrSlug(command.Title, command.Slug); exists == true {
		return nil, domainError.NewDistributionAlreadyExistsError(command.Title)
	}
	if show, _ := service.getShowOutPort.GetShowOrNil(command.ShowId); show == nil {
		return nil, domainError.NewShowNotFoundError(command.ShowId)
	}

	id := uuid.NewString()
	distribution := &model.Distribution{Id: id, ShowId: command.ShowId, Title: command.Title, Slug: command.Slug}
	err := service.saveDistributionOutPort.SaveDistribution(distribution)
	if err != nil {
		return nil, err
	}
	return &onCreateDistribution.CreateDistributionResponse{
		Id:     distribution.Id,
		ShowId: distribution.ShowId,
		Title:  distribution.Title,
		Slug:   distribution.Slug,
	}, nil
}
