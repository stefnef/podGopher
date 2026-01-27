package distribution

import (
	"fmt"
	domainError "podGopher/core/domain/error"
	"podGopher/core/domain/model"
	forGetEpisode "podGopher/core/port/outbound/episode"

	onUpdateDistribution "podGopher/core/port/inbound/distribution"
	forGetSaveDistribution "podGopher/core/port/outbound/distribution"
	forGetShow "podGopher/core/port/outbound/show"
)

type UpdateDistributionService struct {
	getShowOutPort          forGetShow.GetShowPort
	getEpisodeOutPort       forGetEpisode.GetEpisodePort
	getDistributionOutPort  forGetSaveDistribution.GetDistributionPort
	saveDistributionOutPort forGetSaveDistribution.SaveDistributionPort
}

func NewUpdateDistributionService(
	getShowPort forGetShow.GetShowPort,
	getEpisodeOutPort forGetEpisode.GetEpisodePort,
	getDistributionOutPort forGetSaveDistribution.GetDistributionPort,
	saveDistributionPort forGetSaveDistribution.SaveDistributionPort,
) *UpdateDistributionService {
	return &UpdateDistributionService{
		getShowOutPort:          getShowPort,
		getEpisodeOutPort:       getEpisodeOutPort,
		getDistributionOutPort:  getDistributionOutPort,
		saveDistributionOutPort: saveDistributionPort,
	}
}

func (service UpdateDistributionService) UpdateDistribution(command *onUpdateDistribution.UpdateDistributionCommand) error {
	var updatedDistribution *model.Distribution
	var err error

	if updatedDistribution, err = service.getDistributionFromId(command); err != nil {
		return err
	}

	if err = service.setTitleAndSlug(command, updatedDistribution); err != nil {
		return err
	}

	if err = service.setEpisodes(command, updatedDistribution); err != nil {
		return err
	}

	return service.saveDistributionOutPort.UpdateDistribution(updatedDistribution)
}

func (service UpdateDistributionService) setEpisodes(command *onUpdateDistribution.UpdateDistributionCommand, updatedDistribution *model.Distribution) error {
	if command.Episodes != nil {
		for _, episodeId := range *command.Episodes {
			if episode, _ := service.getEpisodeOutPort.GetEpisodeOrNil(episodeId); episode == nil {
				return domainError.NewUpdateError(fmt.Sprintf("Episode with id '%v' does not exist", episodeId))
			} else {
				if episode.ShowId != command.ShowId {
					return domainError.NewUpdateError(fmt.Sprintf("Episode with id '%v' does not belong to show with id '%v'", episodeId, command.ShowId))
				}
			}
		}
		updatedDistribution.Episodes = *command.Episodes
	}
	return nil
}

func (service UpdateDistributionService) getDistributionFromId(command *onUpdateDistribution.UpdateDistributionCommand) (*model.Distribution, error) {
	var distribution *model.Distribution

	if distribution, _ = service.getDistributionOutPort.GetDistributionOrNil(command.DistributionId); distribution == nil || distribution.ShowId != command.ShowId {
		return nil, domainError.NewDistributionNotFoundError(command.DistributionId)
	}

	if show, _ := service.getShowOutPort.GetShowOrNil(command.ShowId); show == nil {
		return nil, domainError.NewShowNotFoundError(command.ShowId)
	}

	updatedDistribution := &model.Distribution{
		Id:       distribution.Id,
		ShowId:   distribution.ShowId,
		Title:    distribution.Title,
		Slug:     distribution.Slug,
		Episodes: distribution.Episodes,
	}
	return updatedDistribution, nil
}

func (service UpdateDistributionService) setTitleAndSlug(command *onUpdateDistribution.UpdateDistributionCommand, updatedDistribution *model.Distribution) error {
	var title string
	var slug string
	if command.Title != nil {
		title = *command.Title
		updatedDistribution.Title = title
	}
	if command.Slug != nil {
		slug = *command.Slug
		updatedDistribution.Slug = slug
	}

	if exists := service.saveDistributionOutPort.ExistsByTitleOrSlug(title, slug); exists == true {
		return domainError.NewUpdateError("Distribution with title or slug already exists")
	}
	return nil
}
