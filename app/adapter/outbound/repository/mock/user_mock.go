package mock

import (
	"podGopher/core/domain/model"
	"testing"
)

type GetAndSaveUserTestAdapter struct {
	t                                  *testing.T
	CalledGet                          int
	CalledSave                         int //TODO different mock variables needed: CalledSave and CalledAssigned
	OnSaveCalledWith                   *model.User
	ReturnsOnGetOrNilUserByUsername    map[string]*model.User
	ReturnsOnGetOrNilUserByUserId      map[string]*model.User
	WithErrorOnSaveUser                error
	ReturnsOnExistsByUsername          map[string]bool
	ReturnsOnExistsByShowIdAndByUserId map[string]bool
}

func (a *GetAndSaveUserTestAdapter) SaveUser(user *model.User) (err error) {
	a.CalledSave++
	a.OnSaveCalledWith = user
	return a.WithErrorOnSaveUser
}

func (a *GetAndSaveUserTestAdapter) GetUserByIdOrNil(id string) (*model.User, error) {
	a.CalledGet++
	user := a.ReturnsOnGetOrNilUserByUserId[id]
	return user, nil
}

func (a *GetAndSaveUserTestAdapter) GetUserByUsernameOrNil(username string) (*model.User, error) {
	a.CalledGet++
	user := a.ReturnsOnGetOrNilUserByUsername[username]
	return user, nil
}

func (a *GetAndSaveUserTestAdapter) ExistsByShowIdAndByUserId(showId string, username string) (exist bool) {
	return a.ReturnsOnExistsByShowIdAndByUserId[showId+username]
}

func (a *GetAndSaveUserTestAdapter) ExistsByUsername(username string) (exist bool) {
	return a.ReturnsOnExistsByUsername[username]
}

func (a *GetAndSaveUserTestAdapter) Init(t *testing.T) {
	a.t = t
	a.CalledGet = 0
	a.CalledSave = 0
	a.OnSaveCalledWith = nil
	a.ReturnsOnGetOrNilUserByUserId = make(map[string]*model.User)
	a.ReturnsOnGetOrNilUserByUsername = make(map[string]*model.User)
	a.WithErrorOnSaveUser = nil
	a.ReturnsOnExistsByUsername = make(map[string]bool)
	a.ReturnsOnExistsByShowIdAndByUserId = make(map[string]bool)
}

func NewGetAndSaveUserTestAdapter(t *testing.T) *GetAndSaveUserTestAdapter {
	adapter := &GetAndSaveUserTestAdapter{}
	adapter.Init(t)
	return adapter
}
