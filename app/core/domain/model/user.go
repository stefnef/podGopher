package model

import "podGopher/core/domain/role"

type User struct {
	Id        string
	Username  string
	IsAdmin   bool
	ShowRoles []domainRole.ShowRole
}
