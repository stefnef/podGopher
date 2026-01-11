package show

import "podGopher/core/domain/model"

type GetShowPort interface {
	GetShowOrNil(id string) (*model.Show, error)
}
