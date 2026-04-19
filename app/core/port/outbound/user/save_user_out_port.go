package user

import "podGopher/core/domain/model"

type SaveUserPort interface {
	SaveUser(user *model.User) (err error)
	ExistsByUsername(username string) (exist bool)
	ExistsByShowIdAndByUserId(showId string, userId string) (exist bool)
}
