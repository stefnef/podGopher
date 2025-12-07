package outbound

import "podGopher/core/domain/model"

type GetDistributionPort interface {
	GetDistributionOrNil(id string) (*model.Distribution, error)
}
