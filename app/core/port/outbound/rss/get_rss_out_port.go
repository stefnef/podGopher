package rss

import "podGopher/core/domain/model"

type GetRSSPort interface {
	GetRSSOrNil(showSlug string, distributionSlug string) (*model.RSS, error)
}
