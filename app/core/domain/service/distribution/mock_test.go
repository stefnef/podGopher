package distribution

import (
	"podGopher/core/domain/model"
	onCreateDistribution "podGopher/core/port/inbound/distribution"
)

type saveAndGetDistributionTestAdapter struct {
	calledGet                       int
	calledSave                      int
	calledUpdate                    int
	calledUpdateWith                *model.Distribution
	onSaveCalledWith                *model.Distribution
	returnsOnExistsByTitleOrSlug    map[string]bool
	withErrorOnSaveDistribution     error
	withErrorOnGetDistributionOrNil error
	returnsOnGetDistributionOrNil   map[string]*model.Distribution
}

type getShowTestAdapter struct {
	called                int
	returnsOnGetOrNilShow map[string]*model.Show
}

type getEpisodeTestAdapter struct {
	called                   int
	returnsOnGetEpisodeOrNil map[string]*model.Episode
}

func newSaveAndGetDistributionTestAdapter() *saveAndGetDistributionTestAdapter {
	adapter := &saveAndGetDistributionTestAdapter{}
	adapter.init()
	return adapter
}

func newTestCreateDistributionCommand(title string) *onCreateDistribution.CreateDistributionCommand {
	distribution := &onCreateDistribution.CreateDistributionCommand{
		ShowId: "test-show-id",
		Title:  title,
		Slug:   "Slug",
	}
	return distribution
}

func (adapter *saveAndGetDistributionTestAdapter) SaveDistribution(distribution *model.Distribution) error {
	adapter.calledSave++
	adapter.onSaveCalledWith = distribution
	return adapter.withErrorOnSaveDistribution
}

func (adapter *saveAndGetDistributionTestAdapter) UpdateDistribution(distribution *model.Distribution) error {
	adapter.calledUpdate++
	adapter.calledUpdateWith = distribution
	return nil
}

func (a *getShowTestAdapter) GetShowOrNil(id string) (*model.Show, error) {
	a.called++
	show := a.returnsOnGetOrNilShow[id]
	return show, nil
}

func (adapter *saveAndGetDistributionTestAdapter) init() {
	adapter.calledGet = 0
	adapter.calledSave = 0
	adapter.calledUpdate = 0
	adapter.calledUpdateWith = nil
	adapter.onSaveCalledWith = nil
	adapter.returnsOnExistsByTitleOrSlug = make(map[string]bool)
	adapter.returnsOnGetDistributionOrNil = make(map[string]*model.Distribution)
	adapter.withErrorOnSaveDistribution = nil
	adapter.withErrorOnGetDistributionOrNil = nil
}

func (adapter *saveAndGetDistributionTestAdapter) everyExistsByTitleReturns(title string, slug string, returnValue bool) {
	adapter.returnsOnExistsByTitleOrSlug[title+slug] = returnValue
}

func (adapter *saveAndGetDistributionTestAdapter) ExistsByTitleOrSlug(title string, slug string) bool {
	return adapter.returnsOnExistsByTitleOrSlug[title+slug]
}

func (adapter *saveAndGetDistributionTestAdapter) GetDistributionOrNil(id string) (*model.Distribution, error) {
	adapter.calledGet++
	return adapter.returnsOnGetDistributionOrNil[id], adapter.withErrorOnGetDistributionOrNil
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

func newGetEpisodeTestAdapter() *getEpisodeTestAdapter {
	adapter := &getEpisodeTestAdapter{}
	adapter.init()
	return adapter
}

func (a *getEpisodeTestAdapter) GetEpisodeOrNil(id string) (*model.Episode, error) {
	a.called++
	return mockGetEpisodeAdapter.returnsOnGetEpisodeOrNil[id], nil
}

func (a *getEpisodeTestAdapter) init() {
	a.called = 0
	a.returnsOnGetEpisodeOrNil = make(map[string]*model.Episode)
}

func initAdapter() {
	mockGetShowAdapter.init()
	mockGetEpisodeAdapter.init()
	mockSaveAndGetDistributionAdapter.init()
}

var mockSaveAndGetDistributionAdapter = newSaveAndGetDistributionTestAdapter()
var mockGetShowAdapter = newGetShowTestAdapter()
var mockGetEpisodeAdapter = newGetEpisodeTestAdapter()

func ptrOfString(value string) *string {
	return &value
}
