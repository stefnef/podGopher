package outbound

import "podGopher/core/domain/model"

type GetShowPort interface {
	GetShowOrNil(Id string) (*model.Show, error)
}
