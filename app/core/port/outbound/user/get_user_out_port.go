package user

import "podGopher/core/domain/model"

type GetUserPort interface {
	GetUserByIdOrNil(username string) (user *model.User, err error)
	GetUserByUsernameOrNil(username string) (user *model.User, err error)
}
