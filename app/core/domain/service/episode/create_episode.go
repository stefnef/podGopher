package episode

import (
	error2 "podGopher/core/domain/error"
	"podGopher/core/domain/model"
	"podGopher/core/port/inbound"
	"podGopher/core/port/outbound"

	"github.com/google/uuid"
)

type CreateEpisodeService struct {
	getShowOutPort     outbound.GetShowPort
	saveEpisodeOutPort outbound.SaveEpisodePort
}

func NewCreateEpisodeService(showRepository outbound.GetShowPort, episodeRepository outbound.SaveEpisodePort) *CreateEpisodeService {
	return &CreateEpisodeService{
		getShowOutPort:     showRepository,
		saveEpisodeOutPort: episodeRepository,
	}
}

func (service CreateEpisodeService) CreateEpisode(command *inbound.CreateEpisodeCommand) (*inbound.CreateEpisodeResponse, error) {
	if exists := service.saveEpisodeOutPort.ExistsByTitle(command.Title); exists != false {
		return nil, error2.NewEpisodeAlreadyExistsError(command.Title)
	}
	if show, _ := service.getShowOutPort.GetShowOrNil(command.ShowId); show == nil {
		return nil, error2.NewShowNotFoundError(command.ShowId)
	}

	id := uuid.NewString()
	episode := &model.Episode{Id: id, ShowId: command.ShowId, Title: command.Title}
	err := service.saveEpisodeOutPort.SaveEpisode(episode)
	if err != nil {
		return nil, err
	}
	return &inbound.CreateEpisodeResponse{
		Id:     episode.Id,
		ShowId: episode.ShowId,
		Title:  episode.Title,
	}, nil
}
