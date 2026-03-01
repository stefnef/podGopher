package user

import "podGopher/core/domain/model"

type CreateUserCommand struct {
	Username string
	ShowId   string
	Role     string
}

type CreateUserResponse struct {
	Id        string
	Username  string
	ShowRoles []model.ShowRole
}

type CreateUserPort interface {
	CreateUser(command *CreateUserCommand) (user *CreateUserResponse, err error)
}
