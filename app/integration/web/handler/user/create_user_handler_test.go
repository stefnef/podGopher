package user

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	domainError "podGopher/core/domain/error"
	"podGopher/core/domain/model"
	"podGopher/core/port/inbound"
	"podGopher/core/port/inbound/user"
	"podGopher/integration/web/auth"
	"podGopher/integration/web/handler"
	"podGopher/integration/web/handler/handlerTestSetup"
	"testing"

	"github.com/stretchr/testify/assert"
)

type createUserTestService struct {
	called              int
	command             *user.CreateUserCommand
	returnsOnCreateUser *user.CreateUserResponse
	failsWith           error
}

func (s *createUserTestService) init() {
	s.called = 0
	s.command = nil
	s.returnsOnCreateUser = nil
	s.failsWith = nil
}

func (s *createUserTestService) CreateUser(command *user.CreateUserCommand) (user *user.CreateUserResponse, err error) {
	s.called++
	s.command = command
	return s.returnsOnCreateUser, s.failsWith
}

var mockCreateUserService = new(createUserTestService)

var createUserHandler = NewCreateUserHandler(inbound.PortMap{
	inbound.CreateUser: mockCreateUserService,
}, auth.NewAdminAuth("user", "password"))

func Test_should_implement_handler_for_create_user(t *testing.T) {
	assert.NotNil(t, createUserHandler)
	assert.Implements(t, (*handler.Handler)(nil), createUserHandler)
}

func Test_should_panic_if_no_port_was_found_on_create_user_handler(t *testing.T) {
	invalidPortMap := inbound.PortMap{
		inbound.PortInvalid: mockCreateUserService,
	}

	assert.Panics(t, func() {
		NewCreateUserHandler(invalidPortMap, auth.NewAdminAuth("user", "password"))
	})
}

func Test_should_return_route_on_create_user(t *testing.T) {
	var route = createUserHandler.GetRoute()

	var expectedRoute = &handler.Route{
		Method: "POST",
		Path:   "/admin/show/:showId/user",
	}

	assert.Equal(t, expectedRoute, route)
}

func Test_should_propagate_error_on_create_user(t *testing.T) {
	defer mockCreateUserService.init()
	var context, _ = handlerTestSetup.GetTestGinContext(t)
	expectedError := errors.New("some error")

	test := struct {
		paramShowId          string
		requestBody          string
		expectedPortResponse error
	}{
		`some-error-id`,
		`{"Role":"Editor", "Username":"some username"}`,
		expectedError,
	}

	mockCreateUserService.failsWith = test.expectedPortResponse

	context.Request = httptest.NewRequest("POST", "/fake", bytes.NewBuffer([]byte(test.requestBody)))

	createUserHandler.Handle(context)

	assert.NotEmpty(t, context.Errors)
	assert.Equal(t, expectedError, (*context.Errors[0]).Err)
}

func Test_should_call_service_on_create_user(t *testing.T) {
	defer mockCreateUserService.init()
	var createUserDto *userResponseDto
	var context, recorder = handlerTestSetup.GetTestGinContext(t)

	test := struct {
		webParameterShowId   string
		webRequestBody       string
		expectedPortCommand  *user.CreateUserCommand
		expectedPortResponse *user.CreateUserResponse
		expectedWebResponse  *userResponseDto
	}{
		`some-show-id`,
		`{"Username":"some username", "Role": "Editor"}`,
		&user.CreateUserCommand{
			ShowId:   "some-show-id",
			Role:     "Editor",
			Username: "some username",
		},
		&user.CreateUserResponse{
			Id:       "some-id",
			Username: "Mocked Username",
			ShowRoles: []model.ShowRole{
				{
					ShowId: "show-id",
					Role:   model.EDITOR,
				},
				{
					ShowId: "other-show-id",
					Role:   model.PRODUCER,
				},
			},
		},
		&userResponseDto{
			Id:       "some-id",
			Username: "Mocked Username",
			ShowRoles: []model.ShowRole{
				{
					ShowId: "show-id",
					Role:   model.EDITOR,
				},
				{
					ShowId: "other-show-id",
					Role:   model.PRODUCER,
				},
			},
		},
	}

	mockCreateUserService.returnsOnCreateUser = test.expectedPortResponse

	context.Request = httptest.NewRequest("POST", "/fake", bytes.NewBuffer([]byte(test.webRequestBody)))
	context.AddParam("showId", test.webParameterShowId)

	createUserHandler.Handle(context)

	var err = json.Unmarshal(recorder.Body.Bytes(), &createUserDto)

	assert.Equal(t, 1, mockCreateUserService.called)
	assert.Equal(t, test.expectedPortCommand, mockCreateUserService.command)
	assert.Nil(t, err)
	assert.Empty(t, context.Errors)
	assert.Equal(t, test.expectedWebResponse, createUserDto)
	assert.Equal(t, http.StatusCreated, recorder.Code)
}

func Test_should_fail_on_missing_credentials(t *testing.T) {
	defer mockCreateUserService.init()
	var context, recorder = handlerTestSetup.GetTestGinContext(t)
	var expectedError = domainError.NewAuthorizationError()

	createUserHandler.Authorize(context)

	assert.NotEmpty(t, context.Errors)
	assert.Equal(t, expectedError, (*context.Errors[0]).Err)
	assert.Equal(t, http.StatusForbidden, recorder.Code)
}

func Test_should_fail_on_invalid_credentials(t *testing.T) {
	defer mockCreateUserService.init()
	var expectedError = domainError.NewAuthorizationError()

	test := []struct {
		name     string
		username string
		password string
	}{
		{
			name:     "invalid password",
			username: "user",
			password: "invalid-password",
		},
		{
			name:     "invalid username",
			username: "invalid-user",
			password: "password",
		},
		{
			name:     "empty credentials",
			username: "",
			password: "",
		},
	}

	for _, tc := range test {
		t.Run(tc.name, func(t *testing.T) {
			var context, recorder = handlerTestSetup.GetTestGinContext(t)
			context.Request.SetBasicAuth(tc.username, tc.password)

			createUserHandler.Authorize(context)

			assert.NotEmpty(t, context.Errors)
			assert.Equal(t, expectedError, (*context.Errors[0]).Err)
			assert.Equal(t, http.StatusForbidden, recorder.Code)
		})
	}

}

func Test_should_pass_on_valid_credentials(t *testing.T) {
	defer mockCreateUserService.init()
	var context, _ = handlerTestSetup.GetTestGinContext(t)
	context.Request.SetBasicAuth("user", "password")

	createUserHandler.Authorize(context)

	assert.Empty(t, context.Errors)
}
