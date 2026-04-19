package user

import (
	domainError "podGopher/core/domain/error"
	"podGopher/core/domain/model"
	"podGopher/core/domain/role"
	onAssignUser "podGopher/core/port/inbound/user"
	"testing"

	"github.com/stretchr/testify/assert"
)

var assignUserService = NewAssignUserService(mockGetShowAdapter, mockGetAndSaveUserAdapter, mockGetAndSaveUserAdapter)

func Test_should_implement_AssignUserInPort(t *testing.T) {
	assert.NotNil(t, assignUserService)
	assert.Implements(t, (*onAssignUser.AssignUserPort)(nil), assignUserService)
}

func Test_should_throw_error_if_user_was_already_assigned(t *testing.T) {
	defer initAdapter(t)

	mockGetAndSaveUserAdapter.ReturnsOnExistsByShowIdAndByUserId["some-shownew-user-id"] = true
	mockGetAndSaveUserAdapter.ReturnsOnGetOrNilUserByUsername["producer"] = &model.User{Id: "mocked-user-id", Username: "producer", ShowRoles: []domainRole.ShowRole{{"some-show", domainRole.PRODUCER}}}

	command := &onAssignUser.AssignUserCommand{
		AssigneeUsername: "producer",
		ShowId:           "some-show",
		UserId:           "new-user-id",
		Role:             "producer",
	}
	result, err := assignUserService.AssignUser(command)

	assert.Nil(t, result)
	assert.NotNil(t, err)
	assert.Equal(t, domainError.NewUserAlreadyAssignedError("some-show", "new-user-id"), err)
	assert.Equal(t, 0, mockGetAndSaveUserAdapter.CalledSave)
}

func Test_should_propagate_errors_from_adapter_on_assign_user(t *testing.T) {
	defer initAdapter(t)

	expectedError := domainError.NewUpdateError("unknown role")
	mockGetAndSaveUserAdapter.WithErrorOnSaveUser = expectedError
	mockGetShowAdapter.ReturnsOnGetOrNilShow["test-show-id"] = &model.Show{Id: "test-show-id"}
	mockGetAndSaveUserAdapter.ReturnsOnGetOrNilUserByUserId["new-user-id"] = &model.User{Id: "mocked-user-id", Username: "new-user"}
	mockGetAndSaveUserAdapter.ReturnsOnGetOrNilUserByUsername["producer"] = &model.User{Id: "mocked-user-id", Username: "producer", ShowRoles: []domainRole.ShowRole{{"test-show-id", domainRole.PRODUCER}}}

	command := &onAssignUser.AssignUserCommand{
		AssigneeUsername: "producer",
		ShowId:           "test-show-id",
		UserId:           "new-user-id",
	}
	result, err := assignUserService.AssignUser(command)

	assert.Nil(t, result)
	assert.NotNil(t, err)
	assert.Equal(t, expectedError, err)
	assert.Equal(t, 0, mockGetAndSaveUserAdapter.CalledSave)
}

func Test_should_throw_error_if_show_does_not_exist_on_assign_user(t *testing.T) {
	defer initAdapter(t)
	command := &onAssignUser.AssignUserCommand{
		AssigneeUsername: "producer",
		ShowId:           "non-existing-show-id",
		UserId:           "any-user-id",
	}

	mockGetShowAdapter.ReturnsOnGetOrNilShow["non-existing-show-id"] = nil
	mockGetAndSaveUserAdapter.ReturnsOnGetOrNilUserByUsername["producer"] = &model.User{Id: "mocked-user-id", Username: "producer", ShowRoles: []domainRole.ShowRole{{"non-existing-show-id", domainRole.PRODUCER}}}

	result, err := assignUserService.AssignUser(command)

	assert.Nil(t, result)
	assert.NotNil(t, err)
	assert.Equal(t, domainError.NewShowNotFoundError("non-existing-show-id"), err)
	assert.Equal(t, 0, mockGetAndSaveUserAdapter.CalledSave)
}

func Test_should_throw_error_if_user_does_not_exist_on_assign_user(t *testing.T) {
	defer initAdapter(t)
	command := &onAssignUser.AssignUserCommand{
		AssigneeUsername: "producer",
		ShowId:           "test-show-id",
		UserId:           "non-existing-user-id",
	}

	mockGetShowAdapter.ReturnsOnGetOrNilShow["test-show-id"] = &model.Show{Id: "test-show-id"}
	mockGetAndSaveUserAdapter.ReturnsOnGetOrNilUserByUsername["non-existing-user-id"] = nil
	mockGetAndSaveUserAdapter.ReturnsOnGetOrNilUserByUsername["producer"] = &model.User{Id: "mocked-user-id", Username: "producer", ShowRoles: []domainRole.ShowRole{{"test-show-id", domainRole.PRODUCER}}}

	result, err := assignUserService.AssignUser(command)

	assert.Nil(t, result)
	assert.NotNil(t, err)
	assert.Equal(t, domainError.NewUserNotFoundError("non-existing-user-id"), err)
	assert.Equal(t, 0, mockGetAndSaveUserAdapter.CalledSave)
}

func Test_should_throw_error_if_role_is_unknown_on_assign_user(t *testing.T) {
	defer initAdapter(t)
	command := &onAssignUser.AssignUserCommand{
		AssigneeUsername: "producer",
		ShowId:           "test-show-id",
		UserId:           "some-user-id",
		Role:             "bad-role",
	}

	mockGetShowAdapter.ReturnsOnGetOrNilShow["test-show-id"] = &model.Show{Id: "test-show-id"}
	mockGetAndSaveUserAdapter.ReturnsOnGetOrNilUserByUserId["some-user-id"] = &model.User{Id: "mocked-user-id", Username: "some-user"}
	mockGetAndSaveUserAdapter.ReturnsOnGetOrNilUserByUsername["producer"] = &model.User{Id: "mocked-user-id", Username: "producer", ShowRoles: []domainRole.ShowRole{{"test-show-id", domainRole.PRODUCER}}}

	result, err := assignUserService.AssignUser(command)

	assert.Nil(t, result)
	assert.NotNil(t, err)
	assert.Equal(t, domainError.NewUpdateError("unknown role"), err)
	assert.Equal(t, 0, mockGetAndSaveUserAdapter.CalledSave)
}

func Test_should_throw_error_if_assignee_is_not_admin_nor_producer_on_assign_user(t *testing.T) {
	defer initAdapter(t)
	command := &onAssignUser.AssignUserCommand{
		AssigneeUsername: "just-editor",
		ShowId:           "test-show-id",
		UserId:           "some-user-id",
		Role:             "EDITOR",
	}

	mockGetShowAdapter.ReturnsOnGetOrNilShow["test-show-id"] = &model.Show{Id: "test-show-id"}
	mockGetAndSaveUserAdapter.ReturnsOnGetOrNilUserByUserId["some-user-id"] = &model.User{Id: "mocked-user-id", Username: "some-user"}
	mockGetAndSaveUserAdapter.ReturnsOnGetOrNilUserByUsername["just-editor"] = &model.User{Id: "editor-id", IsAdmin: false, Username: "some-user",
		ShowRoles: []domainRole.ShowRole{{"test-show-id", domainRole.EDITOR}, {"other-show-id", domainRole.PRODUCER}}}

	result, err := assignUserService.AssignUser(command)

	assert.Nil(t, result)
	assert.NotNil(t, err)
	assert.Equal(t, domainError.NewAuthorizationError(), err)
	assert.Equal(t, 0, mockGetAndSaveUserAdapter.CalledSave)
}

func Test_should_throw_error_if_assignee_does_not_exist_on_assign_user(t *testing.T) {
	defer initAdapter(t)
	command := &onAssignUser.AssignUserCommand{
		AssigneeUsername: "non-existing",
		ShowId:           "test-show-id",
		UserId:           "some-user-id",
		Role:             "EDITOR",
	}

	mockGetShowAdapter.ReturnsOnGetOrNilShow["test-show-id"] = &model.Show{Id: "test-show-id"}
	mockGetAndSaveUserAdapter.ReturnsOnGetOrNilUserByUserId["some-user-id"] = &model.User{Id: "mocked-user-id", Username: "some-user"}

	result, err := assignUserService.AssignUser(command)

	assert.Nil(t, result)
	assert.NotNil(t, err)
	assert.Equal(t, domainError.NewAuthorizationError(), err)
	assert.Equal(t, 0, mockGetAndSaveUserAdapter.CalledSave)
}

func Test_should_save_mapping_if_assignee_is_producer_on_assign_user(t *testing.T) {
	defer initAdapter(t)

	mockGetShowAdapter.ReturnsOnGetOrNilShow["test-show-id"] = &model.Show{Id: "mocked-show-id"}
	mockGetAndSaveUserAdapter.ReturnsOnGetOrNilUserByUsername["producer"] = &model.User{Id: "mocked-user-id", Username: "producer", ShowRoles: []domainRole.ShowRole{{"test-show-id", domainRole.PRODUCER}}}
	mockGetAndSaveUserAdapter.ReturnsOnGetOrNilUserByUserId["new-user-id"] = &model.User{
		Id:       "mocked-user-id",
		Username: "new-user",
		ShowRoles: []domainRole.ShowRole{
			{"other-show-id", domainRole.EDITOR},
		},
	}

	command := &onAssignUser.AssignUserCommand{
		AssigneeUsername: "producer",
		ShowId:           "test-show-id",
		UserId:           "new-user-id",
		Role:             "PRODUCER",
	}

	result, err := assignUserService.AssignUser(command)

	savedUser := mockGetAndSaveUserAdapter.OnSaveCalledWith

	expectedSavedUser := &model.User{
		Id:       savedUser.Id,
		Username: "new-user",
		ShowRoles: []domainRole.ShowRole{
			{"other-show-id", domainRole.EDITOR},
			{"test-show-id", domainRole.PRODUCER},
		},
	}
	assert.NotNil(t, savedUser)
	assert.Equal(t, expectedSavedUser, savedUser)
	assert.Equal(t, 1, mockGetAndSaveUserAdapter.CalledSave)
	assert.NotEmpty(t, savedUser.Id)

	assert.Nil(t, err)
	assert.NotNil(t, result)
	assert.IsType(t, (*onAssignUser.AssignUserResponse)(nil), result)

	expectedAssignedUser := &onAssignUser.AssignUserResponse{Id: savedUser.Id, Username: "new-user", ShowRoles: expectedSavedUser.ShowRoles}
	assert.Equal(t, expectedAssignedUser, result)
}

func Test_should_save_mapping_if_assignee_is_admin_on_assign_user(t *testing.T) {
	defer initAdapter(t)

	mockGetShowAdapter.ReturnsOnGetOrNilShow["test-show-id"] = &model.Show{Id: "mocked-show-id"}
	mockGetAndSaveUserAdapter.ReturnsOnGetOrNilUserByUsername["assignee"] = &model.User{Id: "mocked-user-id", Username: "assignee", IsAdmin: true}
	mockGetAndSaveUserAdapter.ReturnsOnGetOrNilUserByUserId["new-user-id"] = &model.User{
		Id:        "mocked-user-id",
		Username:  "new-user",
		ShowRoles: []domainRole.ShowRole{},
	}

	command := &onAssignUser.AssignUserCommand{
		AssigneeUsername: "assignee",
		ShowId:           "test-show-id",
		UserId:           "new-user-id",
		Role:             "PRODUCER",
	}

	result, err := assignUserService.AssignUser(command)

	assert.Equal(t, 1, mockGetAndSaveUserAdapter.CalledSave)
	assert.Nil(t, err)
	assert.NotNil(t, result)
}
