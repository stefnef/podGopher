package mock

import (
	"podGopher/core/domain/model"
	"testing"
)

type GetAndSaveUserTestAdapter struct {
	t                         *testing.T
	CalledGet                 int
	CalledSave                int
	OnSaveCalledWith          *model.User
	ReturnsOnGetOrNilUser     map[string]*model.User
	WithErrorOnSaveUser       error
	ReturnsOnExistsByUsername map[string]bool
}

func (a *GetAndSaveUserTestAdapter) SaveUser(user *model.User) (err error) {
	a.CalledSave++
	a.OnSaveCalledWith = user
	return a.WithErrorOnSaveUser
}

func (a *GetAndSaveUserTestAdapter) GetUserOrNil(id string) (*model.User, error) {
	a.CalledGet++
	user := a.ReturnsOnGetOrNilUser[id]
	return user, nil
}

func (a *GetAndSaveUserTestAdapter) ExistsByUsername(showId string, username string) (exist bool) {
	return a.ReturnsOnExistsByUsername[showId+username]
}

func (a *GetAndSaveUserTestAdapter) Init(t *testing.T) {
	a.t = t
	a.CalledGet = 0
	a.CalledSave = 0
	a.OnSaveCalledWith = nil
	a.ReturnsOnGetOrNilUser = make(map[string]*model.User)
	a.WithErrorOnSaveUser = nil
	a.ReturnsOnExistsByUsername = make(map[string]bool)
}

func NewGetAndSaveUserTestAdapter(t *testing.T) *GetAndSaveUserTestAdapter {
	adapter := &GetAndSaveUserTestAdapter{}
	adapter.Init(t)
	return adapter
}
