package episode

import "podGopher/core/domain/model"

type SaveEpisodePort interface {
	SaveEpisode(episode *model.Episode) (err error)
	ExistsByTitle(title string) (exist bool)
}
