package episode

import (
	domainError "podGopher/core/domain/error"
	"podGopher/core/domain/model"
	onGetEpisode "podGopher/core/port/inbound/episode"
	forGetEpisode "podGopher/core/port/outbound/episode"
	forGetShow "podGopher/core/port/outbound/show"
)

type GetEpisodeService struct {
	getShowOutPort    forGetShow.GetShowPort
	getEpisodeOutPort forGetEpisode.GetEpisodePort
}

func NewGetEpisodeService(showRepository forGetShow.GetShowPort, episodeRepository forGetEpisode.GetEpisodePort) *GetEpisodeService {
	return &GetEpisodeService{
		getShowOutPort:    showRepository,
		getEpisodeOutPort: episodeRepository,
	}
}

func (service *GetEpisodeService) GetEpisode(command *onGetEpisode.GetEpisodeCommand) (episode *onGetEpisode.GetEpisodeResponse, err error) {
	if show, _ := service.getShowOutPort.GetShowOrNil(command.ShowId); show == nil {
		return nil, domainError.NewShowNotFoundError(command.ShowId)
	}

	var foundEpisode *model.Episode
	if foundEpisode, err = service.getEpisodeOutPort.GetEpisodeOrNil(command.EpisodeId); err != nil {
		return nil, err
	}

	if foundEpisode == nil {
		return nil, domainError.NewEpisodeNotFoundError(command.EpisodeId)
	}

	return &onGetEpisode.GetEpisodeResponse{
		Id:     foundEpisode.Id,
		ShowId: foundEpisode.ShowId,
		Title:  foundEpisode.Title,
	}, nil
}
