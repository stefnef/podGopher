package distribution

import "podGopher/core/domain/model"

type GetDistributionPort interface {
	GetDistributionOrNil(id string) (*model.Distribution, error)
}
