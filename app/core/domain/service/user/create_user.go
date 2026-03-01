package user

import (
	domainError "podGopher/core/domain/error"
	"podGopher/core/domain/model"
	onCreateUser "podGopher/core/port/inbound/user"
	forGetShow "podGopher/core/port/outbound/show"
	forGetSaveUser "podGopher/core/port/outbound/user"

	"github.com/google/uuid"
)

type CreateUserService struct {
	getShowOutPort  forGetShow.GetShowPort
	getSaveUserPort forGetSaveUser.SaveUserPort
}

func NewCreateUserService(showRepository forGetShow.GetShowPort, userRepository forGetSaveUser.SaveUserPort) *CreateUserService {
	return &CreateUserService{
		getShowOutPort:  showRepository,
		getSaveUserPort: userRepository,
	}
}

func (service CreateUserService) CreateUser(command *onCreateUser.CreateUserCommand) (*onCreateUser.CreateUserResponse, error) {
	if exists := service.getSaveUserPort.ExistsByUsername(command.ShowId, command.Username); exists != false {
		return nil, domainError.NewUserAlreadyExistsError(command.ShowId, command.Username)
	}

	if show, _ := service.getShowOutPort.GetShowOrNil(command.ShowId); show == nil {
		return nil, domainError.NewShowNotFoundError(command.ShowId)
	}

	id := uuid.NewString()
	user := &model.User{
		Id:       id,
		Username: command.Username,
		ShowRoles: []model.ShowRole{
			{ShowId: command.ShowId, Role: model.PRODUCER},
		},
	}

	if err := service.getSaveUserPort.SaveUser(user); err != nil {
		return nil, err
	}

	return &onCreateUser.CreateUserResponse{
		Id:        user.Id,
		Username:  user.Username,
		ShowRoles: user.ShowRoles,
	}, nil

}
