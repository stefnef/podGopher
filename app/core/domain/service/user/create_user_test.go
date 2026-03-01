package user

import (
	"errors"
	domainError "podGopher/core/domain/error"
	"podGopher/core/domain/model"
	onCreateUser "podGopher/core/port/inbound/user"
	"testing"

	"github.com/stretchr/testify/assert"
)

var createUserService = NewCreateUserService(mockGetShowAdapter, mockGetAndSaveUserAdapter)

func Test_should_implement_CreateUserInPort(t *testing.T) {
	assert.NotNil(t, createUserService)
	assert.Implements(t, (*onCreateUser.CreateUserPort)(nil), createUserService)
}

func Test_should_throw_error_if_episode_with_name_already_exists(t *testing.T) {
	defer initAdapter(t)

	mockGetAndSaveUserAdapter.ReturnsOnExistsByUsername["some-shownew-user"] = true

	command := &onCreateUser.CreateUserCommand{
		ShowId:   "some-show",
		Username: "new-user",
		Role:     "producer",
	}
	result, err := createUserService.CreateUser(command)

	assert.Nil(t, result)
	assert.NotNil(t, err)
	assert.Equal(t, domainError.NewUserAlreadyExistsError("some-show", "new-user"), err)
	assert.Equal(t, 0, mockGetAndSaveUserAdapter.CalledSave)
}

func Test_should_propagate_errors_from_adapter_on_create_user(t *testing.T) {
	defer initAdapter(t)

	expectedError := errors.New("some error")
	mockGetAndSaveUserAdapter.WithErrorOnSaveUser = expectedError
	mockGetShowAdapter.ReturnsOnGetOrNilShow["test-show-id"] = &model.Show{Id: "test-show-id"}

	command := &onCreateUser.CreateUserCommand{
		ShowId:   "test-show-id",
		Username: "new-user",
	}
	result, err := createUserService.CreateUser(command)

	assert.Nil(t, result)
	assert.NotNil(t, err)
	assert.Equal(t, expectedError, err)
	assert.Equal(t, 1, mockGetAndSaveUserAdapter.CalledSave)
}

func Test_should_throw_error_if_show_does_not_exist_on_save_user(t *testing.T) {
	defer initAdapter(t)
	command := &onCreateUser.CreateUserCommand{
		ShowId:   "non-existing-show-id",
		Username: "any-user",
	}

	mockGetShowAdapter.ReturnsOnGetOrNilShow["non-existing-show-id"] = nil

	result, err := createUserService.CreateUser(command)

	assert.Nil(t, result)
	assert.NotNil(t, err)
	assert.Equal(t, domainError.NewShowNotFoundError("non-existing-show-id"), err)
	assert.Equal(t, 0, mockGetAndSaveUserAdapter.CalledSave)
}

func Test_should_save_a_new_user(t *testing.T) {
	defer initAdapter(t)

	mockGetShowAdapter.ReturnsOnGetOrNilShow["test-show-id"] = &model.Show{Id: "mocked-show-id"}

	command := &onCreateUser.CreateUserCommand{
		ShowId:   "test-show-id",
		Username: "new-user",
		Role:     "producer",
	}

	result, err := createUserService.CreateUser(command)

	savedUser := mockGetAndSaveUserAdapter.OnSaveCalledWith

	expectedSavedUser := &model.User{
		Id:       savedUser.Id,
		Username: "new-user",
		ShowRoles: []model.ShowRole{
			{ShowId: "test-show-id", Role: model.PRODUCER},
		},
	}
	assert.NotNil(t, savedUser)
	assert.Equal(t, expectedSavedUser, savedUser)
	assert.Equal(t, 1, mockGetAndSaveUserAdapter.CalledSave)
	assert.NotEmpty(t, savedUser.Id)

	assert.Nil(t, err)
	assert.NotNil(t, result)
	assert.IsType(t, (*onCreateUser.CreateUserResponse)(nil), result)

	expectedCreatedUser := &onCreateUser.CreateUserResponse{Id: savedUser.Id, Username: "new-user", ShowRoles: expectedSavedUser.ShowRoles}
	assert.Equal(t, expectedCreatedUser, result)
}
