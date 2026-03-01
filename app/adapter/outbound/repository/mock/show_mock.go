package mock

import (
	"podGopher/core/domain/model"
	"testing"
)

type GetShowTestAdapter struct {
	t                     *testing.T
	Called                int
	ReturnsOnGetOrNilShow map[string]*model.Show
}

func (a *GetShowTestAdapter) GetAllShows() ([]*model.Show, error) {
	panic("Don't use me")
}

func (a *GetShowTestAdapter) GetShowOrNil(id string) (*model.Show, error) {
	a.Called++
	show := a.ReturnsOnGetOrNilShow[id]
	return show, nil
}

func (a *GetShowTestAdapter) Init(t *testing.T) {
	a.t = t
	a.Called = 0
	a.ReturnsOnGetOrNilShow = make(map[string]*model.Show)
}

func NewGetShowTestAdapter(t *testing.T) *GetShowTestAdapter {
	adapter := &GetShowTestAdapter{}
	adapter.Init(t)
	return adapter
}
