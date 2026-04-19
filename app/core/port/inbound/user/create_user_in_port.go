package user

import (
	"podGopher/core/domain/role"
)

type CreateUserCommand struct {
	Username string
	Email    string
	Password string
	IsAdmin  bool
}

type CreateUserResponse struct {
	Id        string
	Username  string
	ShowRoles []domainRole.ShowRole
}

type CreateUserPort interface {
	CreateUser(command *CreateUserCommand) (user *CreateUserResponse, err error)
}
