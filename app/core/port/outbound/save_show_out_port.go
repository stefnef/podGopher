package outbound

import "podGopher/core/domain/model"

type SaveShowPort interface {
	SaveShow(show *model.Show) (err error)
	ExistsByTitleOrSlug(title string, slug string) bool
}
