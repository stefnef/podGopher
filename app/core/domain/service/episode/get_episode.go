package episode

import (
	error2 "podGopher/core/domain/error"
	"podGopher/core/domain/model"
	"podGopher/core/port/inbound"
	"podGopher/core/port/outbound"
)

type GetEpisodeService struct {
	getShowOutPort    outbound.GetShowPort
	getEpisodeOutPort outbound.GetEpisodePort
}

func NewGetEpisodeService(showRepository outbound.GetShowPort, episodeRepository outbound.GetEpisodePort) *GetEpisodeService {
	return &GetEpisodeService{
		getShowOutPort:    showRepository,
		getEpisodeOutPort: episodeRepository,
	}
}

func (service *GetEpisodeService) GetEpisode(command *inbound.GetEpisodeCommand) (episode *inbound.GetEpisodeResponse, err error) {
	if show, _ := service.getShowOutPort.GetShowOrNil(command.ShowId); show == nil {
		return nil, error2.NewShowNotFoundError(command.ShowId)
	}

	var foundEpisode *model.Episode
	if foundEpisode, err = service.getEpisodeOutPort.GetEpisodeOrNil(command.EpisodeId); err != nil {
		return nil, err
	}

	if foundEpisode == nil {
		return nil, error2.NewEpisodeNotFoundError(command.EpisodeId)
	}

	return &inbound.GetEpisodeResponse{
		Id:     foundEpisode.Id,
		ShowId: foundEpisode.ShowId,
		Title:  foundEpisode.Title,
	}, nil
}
