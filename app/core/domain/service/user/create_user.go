package user

import (
	domainError "podGopher/core/domain/error"
	"podGopher/core/domain/model"
	"podGopher/core/domain/role"
	onCreateUser "podGopher/core/port/inbound/user"
	forGetSaveUser "podGopher/core/port/outbound/user"

	"github.com/google/uuid"
)

type CreateUserService struct {
	getSaveUserPort           forGetSaveUser.SaveUserPort
	createUserCredentialsPort forGetSaveUser.CreateUserCredentialsPort
}

func NewCreateUserService(userRepository forGetSaveUser.SaveUserPort, userCredentialService forGetSaveUser.CreateUserCredentialsPort) *CreateUserService {
	return &CreateUserService{
		getSaveUserPort:           userRepository,
		createUserCredentialsPort: userCredentialService,
	}
}

func (service CreateUserService) CreateUser(command *onCreateUser.CreateUserCommand) (*onCreateUser.CreateUserResponse, error) {
	if exists := service.getSaveUserPort.ExistsByUsername(command.Username); exists != false {
		return nil, domainError.NewUserAlreadyExistsError(command.Username)
	}

	id := uuid.NewString()
	user := &model.User{
		Id:        id,
		Username:  command.Username,
		IsAdmin:   command.IsAdmin,
		ShowRoles: []domainRole.ShowRole{},
	}

	if err := service.createUserCredentialsPort.CreateUserCredentials(user, command.Email, command.Password); err != nil {
		return nil, err
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
