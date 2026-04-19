package user

import (
	"podGopher/core/domain/role"
)

type AssignUserCommand struct {
	AssigneeUsername string
	UserId           string
	ShowId           string
	Role             string
}

type AssignUserResponse struct {
	Id        string
	Username  string
	ShowRoles []domainRole.ShowRole
}

type AssignUserPort interface {
	AssignUser(command *AssignUserCommand) (user *AssignUserResponse, err error)
}
