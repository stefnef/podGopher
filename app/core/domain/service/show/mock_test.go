package show

import (
	"podGopher/core/domain/model"
	"podGopher/core/port/inbound"
)

type saveAndGetShowTestAdapter struct {
	calledSave                   int
	onSave                       map[string]*model.Show
	returnsOnExistsByTitleOrSlug map[string]bool
	withErrorOnSaveShow          error
}

func newSaveAndGetShowTestAdapter() *saveAndGetShowTestAdapter {
	adapter := &saveAndGetShowTestAdapter{}
	adapter.init()
	return adapter
}

func newTestCreateShowCommand(title string) *inbound.CreateShowCommand {
	show := &inbound.CreateShowCommand{
		Title: title,
		Slug:  title + "-Slug",
	}
	return show
}

func (adapter *saveAndGetShowTestAdapter) SaveShow(show *model.Show) error {
	adapter.calledSave++
	adapter.onSave["show"] = show
	return adapter.withErrorOnSaveShow
}

func (adapter *saveAndGetShowTestAdapter) init() {
	adapter.calledSave = 0
	adapter.onSave = make(map[string]*model.Show)
	adapter.returnsOnExistsByTitleOrSlug = make(map[string]bool)
	adapter.withErrorOnSaveShow = nil
}

func (adapter *saveAndGetShowTestAdapter) everyExistsByTitleOrSlugReturns(title string, slug string, returnValue bool) {
	adapter.returnsOnExistsByTitleOrSlug[title+slug] = returnValue
}

func (adapter *saveAndGetShowTestAdapter) ExistsByTitleOrSlug(title string, slug string) bool {
	return adapter.returnsOnExistsByTitleOrSlug[title+slug]
}

type getShowTestAdapter struct {
	called                  int
	returnsOnGetOrNilShow   map[string]*model.Show
	withErrorOnGetOrNilShow error
}

func (a *getShowTestAdapter) GetShowOrNil(Id string) (*model.Show, error) {
	a.called++
	show := a.returnsOnGetOrNilShow[Id]
	return show, a.withErrorOnGetOrNilShow
}

func initAdapter() {
	mockGetShowAdapter.init()
	mockSaveAndGetShowAdapter.init()
}

var mockGetShowAdapter = newGetShowTestAdapter()

var mockSaveAndGetShowAdapter = newSaveAndGetShowTestAdapter()
