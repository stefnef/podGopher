package user

import "podGopher/core/domain/model"

type CreateUserCredentialsPort interface {
	CreateUserCredentials(user *model.User, email string, password string) (err error)
}
