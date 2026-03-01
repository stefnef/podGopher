package user

import "podGopher/core/domain/model"

type SaveUserPort interface {
	SaveUser(user *model.User) (err error)
	ExistsByUsername(showId string, username string) (exist bool)
}
