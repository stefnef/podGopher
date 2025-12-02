package outbound

import "podGopher/core/domain/model"

type GetEpisodePort interface {
	GetEpisodeOrNil(id string) (*model.Episode, error)
}
