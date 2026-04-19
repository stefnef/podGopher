package user

import (
	"errors"
	domainError "podGopher/core/domain/error"
	"podGopher/core/domain/model"
	"podGopher/core/domain/role"
	onCreateUser "podGopher/core/port/inbound/user"
	"testing"

	"github.com/stretchr/testify/assert"
)

var createUserService = NewCreateUserService(mockGetAndSaveUserAdapter, mockCreateUserCredentialsAdapter)

func Test_should_implement_CreateUserInPort(t *testing.T) {
	assert.NotNil(t, createUserService)
	assert.Implements(t, (*onCreateUser.CreateUserPort)(nil), createUserService)
}

func Test_should_throw_error_if_user_with_username_already_exists(t *testing.T) {
	defer initAdapter(t)

	mockGetAndSaveUserAdapter.ReturnsOnExistsByUsername["new-user"] = true

	command := &onCreateUser.CreateUserCommand{
		Username: "new-user",
	}
	result, err := createUserService.CreateUser(command)

	assert.Nil(t, result)
	assert.NotNil(t, err)
	assert.Equal(t, domainError.NewUserAlreadyExistsError("new-user"), err)
	assert.Equal(t, 0, mockGetAndSaveUserAdapter.CalledSave)
}

func Test_should_propagate_errors_from_save_adapter_on_create_user(t *testing.T) {
	defer initAdapter(t)

	expectedError := errors.New("some error")
	mockGetAndSaveUserAdapter.WithErrorOnSaveUser = expectedError

	command := &onCreateUser.CreateUserCommand{
		Username: "new-user",
	}
	result, err := createUserService.CreateUser(command)

	assert.Nil(t, result)
	assert.NotNil(t, err)
	assert.Equal(t, expectedError, err)
	assert.Equal(t, 1, mockGetAndSaveUserAdapter.CalledSave)
}

func Test_should_not_create_if_user_management_adapter_fails_on_create_user(t *testing.T) {
	defer initAdapter(t)

	expectedError := errors.New("some error")
	mockCreateUserCredentialsAdapter.WithErrorOnCreateUserCredentials = expectedError

	command := &onCreateUser.CreateUserCommand{
		Username: "new-user",
	}
	result, err := createUserService.CreateUser(command)

	assert.Nil(t, result)
	assert.NotNil(t, err)
	assert.Equal(t, expectedError, err)
	assert.Equal(t, 0, mockGetAndSaveUserAdapter.CalledSave)
}

func Test_should_save_a_new_user(t *testing.T) {
	defer initAdapter(t)

	command := &onCreateUser.CreateUserCommand{
		Username: "new-user",
		Email:    "some-email",
		Password: "some-password",
		IsAdmin:  false,
	}

	result, err := createUserService.CreateUser(command)

	savedUser := mockGetAndSaveUserAdapter.OnSaveCalledWith

	expectedSavedUser := &model.User{
		Id:        savedUser.Id,
		Username:  "new-user",
		IsAdmin:   false,
		ShowRoles: []domainRole.ShowRole{},
	}
	assert.NotNil(t, savedUser)
	assert.Equal(t, expectedSavedUser, savedUser)
	assert.Equal(t, 1, mockGetAndSaveUserAdapter.CalledSave)
	assert.Equal(t, 1, mockCreateUserCredentialsAdapter.CalledCreate)
	assert.NotEmpty(t, savedUser.Id)

	assert.Nil(t, err)
	assert.NotNil(t, result)
	assert.IsType(t, (*onCreateUser.CreateUserResponse)(nil), result)

	expectedCreatedUser := &onCreateUser.CreateUserResponse{Id: savedUser.Id, Username: "new-user", ShowRoles: []domainRole.ShowRole{}}
	assert.Equal(t, expectedCreatedUser, result)

	assert.Equal(t, expectedSavedUser, mockCreateUserCredentialsAdapter.OnCreateCalledWithUser)
	assert.Equal(t, "some-email", mockCreateUserCredentialsAdapter.OnCreateCalledWithEmail)
	assert.Equal(t, "some-password", mockCreateUserCredentialsAdapter.OnCreateCalledWithPassword)
}

func Test_should_save_is_admin(t *testing.T) {
	defer initAdapter(t)

	command := &onCreateUser.CreateUserCommand{
		Username: "new-user",
		Email:    "some-other-email",
		Password: "some-other-password",
		IsAdmin:  true,
	}

	_, err := createUserService.CreateUser(command)

	assert.Nil(t, err)
	assert.True(t, mockCreateUserCredentialsAdapter.OnCreateCalledWithUser.IsAdmin)
	assert.Equal(t, "some-other-email", mockCreateUserCredentialsAdapter.OnCreateCalledWithEmail)
	assert.Equal(t, "some-other-password", mockCreateUserCredentialsAdapter.OnCreateCalledWithPassword)
}
