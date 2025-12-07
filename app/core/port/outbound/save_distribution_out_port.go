package outbound

import "podGopher/core/domain/model"

type SaveDistributionPort interface {
	SaveDistribution(distribution *model.Distribution) (err error)
	ExistsByTitleOrSlug(title string, slug string) (exist bool)
}
