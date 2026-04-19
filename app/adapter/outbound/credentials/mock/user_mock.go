package mock

import (
	"podGopher/core/domain/model"
	"testing"
)

type CreateUserCredentialsTestAdapter struct {
	t                                *testing.T
	CalledCreate                     int
	OnCreateCalledWithUser           *model.User
	OnCreateCalledWithEmail          string
	OnCreateCalledWithPassword       string
	WithErrorOnCreateUserCredentials error
}

func (a *CreateUserCredentialsTestAdapter) CreateUserCredentials(user *model.User, email string, password string) (err error) {
	a.CalledCreate++
	a.OnCreateCalledWithUser = user
	a.OnCreateCalledWithEmail = email
	a.OnCreateCalledWithPassword = password
	return a.WithErrorOnCreateUserCredentials
}

func (a *CreateUserCredentialsTestAdapter) Init(t *testing.T) {
	a.t = t
	a.CalledCreate = 0
	a.OnCreateCalledWithUser = nil
	a.OnCreateCalledWithEmail = ""
	a.OnCreateCalledWithPassword = ""
	a.WithErrorOnCreateUserCredentials = nil
}

func NewCreateUserTestAdapter(t *testing.T) *CreateUserCredentialsTestAdapter {
	adapter := &CreateUserCredentialsTestAdapter{}
	adapter.Init(t)
	return adapter
}
