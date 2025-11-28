package episode

import (
	"podGopher/core/domain/model"
	"podGopher/core/port/inbound"
)

type saveAndGetEpisodeTestAdapter struct {
	calledGet                  int
	calledSave                 int
	onSaveCalledWith           *model.Episode
	returnsOnExistsByTitle     map[string]bool
	withErrorOnSaveEpisode     error
	withErrorOnGetEpisodeOrNil error
	returnsOnGetEpisodeOrNil   map[string]*model.Episode
}

type getShowTestAdapter struct {
	called                int
	returnsOnGetOrNilShow map[string]*model.Show
}

func newSaveAndGetEpisodeTestAdapter() *saveAndGetEpisodeTestAdapter {
	adapter := &saveAndGetEpisodeTestAdapter{}
	adapter.init()
	return adapter
}

func newTestCreateEpisodeCommand(title string) *inbound.CreateEpisodeCommand {
	episode := &inbound.CreateEpisodeCommand{
		ShowId: "test-show-id",
		Title:  title,
	}
	return episode
}

func (adapter *saveAndGetEpisodeTestAdapter) SaveEpisode(episode *model.Episode) error {
	adapter.calledSave++
	adapter.onSaveCalledWith = episode
	return adapter.withErrorOnSaveEpisode
}

func (a *getShowTestAdapter) GetShowOrNil(id string) (*model.Show, error) {
	a.called++
	show := a.returnsOnGetOrNilShow[id]
	return show, nil
}

func (adapter *saveAndGetEpisodeTestAdapter) init() {
	adapter.calledGet = 0
	adapter.calledSave = 0
	adapter.onSaveCalledWith = nil
	adapter.returnsOnExistsByTitle = make(map[string]bool)
	adapter.returnsOnGetEpisodeOrNil = make(map[string]*model.Episode)
	adapter.withErrorOnSaveEpisode = nil
	adapter.withErrorOnGetEpisodeOrNil = nil
}

func (adapter *saveAndGetEpisodeTestAdapter) everyExistsByTitleReturns(title string, returnValue bool) {
	adapter.returnsOnExistsByTitle[title] = returnValue
}

func (adapter *saveAndGetEpisodeTestAdapter) ExistsByTitle(title string) bool {
	return adapter.returnsOnExistsByTitle[title]
}

func (adapter *saveAndGetEpisodeTestAdapter) GetEpisodeOrNil(id string) (*model.Episode, error) {
	adapter.calledGet++
	return adapter.returnsOnGetEpisodeOrNil[id], adapter.withErrorOnGetEpisodeOrNil
}

func (a *getShowTestAdapter) init() {
	a.called = 0
	a.returnsOnGetOrNilShow = make(map[string]*model.Show)
}

func newGetShowTestAdapter() *getShowTestAdapter {
	adapter := &getShowTestAdapter{}
	adapter.init()
	return adapter
}

func initAdapter() {
	mockGetShowAdapter.init()
	mockSaveAndGetEpisodeAdapter.init()
}

var mockSaveAndGetEpisodeAdapter = newSaveAndGetEpisodeTestAdapter()
var mockGetShowAdapter = newGetShowTestAdapter()
