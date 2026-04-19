package user

import (
	domainError "podGopher/core/domain/error"
	"podGopher/core/domain/model"
	domainRole "podGopher/core/domain/role"
	onAssignUser "podGopher/core/port/inbound/user"
	forGetShow "podGopher/core/port/outbound/show"
	forGetAndSaveUser "podGopher/core/port/outbound/user"
	"slices"
)

type AssignUserService struct {
	getShowOutPort forGetShow.GetShowPort
	getUserPort    forGetAndSaveUser.GetUserPort
	saveUserPort   forGetAndSaveUser.SaveUserPort
}

func NewAssignUserService(showRepository forGetShow.GetShowPort, saveUserRepository forGetAndSaveUser.SaveUserPort, getUserRepository forGetAndSaveUser.GetUserPort) *AssignUserService {
	return &AssignUserService{
		getShowOutPort: showRepository,
		getUserPort:    getUserRepository,
		saveUserPort:   saveUserRepository,
	}
}

func (service AssignUserService) AssignUser(command *onAssignUser.AssignUserCommand) (*onAssignUser.AssignUserResponse, error) {
	user, err := service.validateAndGetUser(command)
	if err != nil {
		return nil, err
	}

	if err = service.saveUserPort.SaveUser(user); err != nil {
		return nil, err
	}

	return &onAssignUser.AssignUserResponse{
		Id:        user.Id,
		Username:  user.Username,
		ShowRoles: user.ShowRoles,
	}, nil

}

func (service AssignUserService) validateAndGetUser(command *onAssignUser.AssignUserCommand) (*model.User, error) {
	if service.isAssigneePermitted(command.AssigneeUsername, command.ShowId) == false {
		return nil, domainError.NewAuthorizationError()
	}

	if err := service.validateShowAndUsername(command.ShowId, command.UserId); err != nil {
		return nil, err
	}

	return service.getUserAndApplyNewRole(command.UserId, command.ShowId, command.Role)
}

func (service AssignUserService) getUserAndApplyNewRole(userId, showId, roleValue string) (*model.User, error) {
	var user *model.User
	if user, _ = service.getUserPort.GetUserByIdOrNil(userId); user == nil {
		return nil, domainError.NewUserNotFoundError(userId)
	}

	role := domainRole.ValueToRole(roleValue)
	if role == domainRole.UNKNOWN {
		return nil, domainError.NewUpdateError("unknown role")
	}
	user.ShowRoles = append(user.ShowRoles, domainRole.ShowRole{ShowId: showId, Role: role})

	return user, nil
}

func (service AssignUserService) validateShowAndUsername(showId string, username string) error {
	if exists := service.saveUserPort.ExistsByShowIdAndByUserId(showId, username); exists != false {
		return domainError.NewUserAlreadyAssignedError(showId, username)
	}

	if show, _ := service.getShowOutPort.GetShowOrNil(showId); show == nil {
		return domainError.NewShowNotFoundError(showId)
	}
	return nil
}

func (service AssignUserService) isAssigneePermitted(assigneeUsername string, showId string) bool {
	var assignee *model.User
	if assignee, _ = service.getUserPort.GetUserByUsernameOrNil(assigneeUsername); assignee == nil {
		return false
	}
	if !assignee.IsAdmin && !slices.Contains(assignee.ShowRoles, domainRole.ShowRole{ShowId: showId, Role: domainRole.PRODUCER}) {
		return false
	}
	return true
}
