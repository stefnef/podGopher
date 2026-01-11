package episode

import (
	domainError "podGopher/core/domain/error"
	"podGopher/core/domain/model"
	onCreateEpisode "podGopher/core/port/inbound/episode"
	forSaveEpisode "podGopher/core/port/outbound/episode"
	forGetShow "podGopher/core/port/outbound/show"

	"github.com/google/uuid"
)

type CreateEpisodeService struct {
	getShowOutPort     forGetShow.GetShowPort
	saveEpisodeOutPort forSaveEpisode.SaveEpisodePort
}

func NewCreateEpisodeService(showRepository forGetShow.GetShowPort, episodeRepository forSaveEpisode.SaveEpisodePort) *CreateEpisodeService {
	return &CreateEpisodeService{
		getShowOutPort:     showRepository,
		saveEpisodeOutPort: episodeRepository,
	}
}

func (service CreateEpisodeService) CreateEpisode(command *onCreateEpisode.CreateEpisodeCommand) (*onCreateEpisode.CreateEpisodeResponse, error) {
	if exists := service.saveEpisodeOutPort.ExistsByTitle(command.Title); exists != false {
		return nil, domainError.NewEpisodeAlreadyExistsError(command.Title)
	}
	if show, _ := service.getShowOutPort.GetShowOrNil(command.ShowId); show == nil {
		return nil, domainError.NewShowNotFoundError(command.ShowId)
	}

	id := uuid.NewString()
	episode := &model.Episode{Id: id, ShowId: command.ShowId, Title: command.Title}
	err := service.saveEpisodeOutPort.SaveEpisode(episode)
	if err != nil {
		return nil, err
	}
	return &onCreateEpisode.CreateEpisodeResponse{
		Id:     episode.Id,
		ShowId: episode.ShowId,
		Title:  episode.Title,
	}, nil
}
