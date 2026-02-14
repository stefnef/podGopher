package rss

import "podGopher/core/domain/model"

type getRSSTestAdapter struct {
	called                 int
	returnsOnGetOrNilError map[string]error
	returnsOnGetOrNilRSS   map[string]*model.RSS
}

func (a *getRSSTestAdapter) GetRSSOrNil(showSlug string, distributionSlug string) (*model.RSS, error) {
	a.called++
	return a.returnsOnGetOrNilRSS[showSlug+distributionSlug], a.returnsOnGetOrNilError[showSlug+distributionSlug]
}

func newGetRSSTestAdapter() *getRSSTestAdapter {
	adapter := &getRSSTestAdapter{}
	adapter.init()
	return adapter
}

func (a *getRSSTestAdapter) init() {
	a.called = 0
	a.returnsOnGetOrNilError = make(map[string]error)
	a.returnsOnGetOrNilRSS = make(map[string]*model.RSS)
}

var mockGetRSSAdapter = newGetRSSTestAdapter()

func initAdapter() {
	mockGetRSSAdapter.init()
}
