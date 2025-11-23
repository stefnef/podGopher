package show

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"podGopher/core/port/inbound"
	"podGopher/integration/web/handler"
	"podGopher/integration/web/handler/handlerTestSetup"
	"testing"

	"github.com/stretchr/testify/assert"
)

type createShowTestService struct {
	called              int
	command             *inbound.CreateShowCommand
	returnsOnCreateShow *inbound.CreateShowResponse
	failsWith           error
}

func (s *createShowTestService) init() {
	s.called = 0
	s.command = nil
	s.returnsOnCreateShow = nil
	s.failsWith = nil
}

func (s *createShowTestService) CreateShow(command *inbound.CreateShowCommand) (show *inbound.CreateShowResponse, err error) {
	s.called++
	s.command = command
	return s.returnsOnCreateShow, s.failsWith
}

var mockCreateShowService = new(createShowTestService)
var createShowHandler = NewCreateShowHandler(inbound.PortMap{
	inbound.CreateShow: mockCreateShowService,
})

func Test_should_implement_handler_for_create_show(t *testing.T) {
	assert.NotNil(t, createShowHandler)
	assert.Implements(t, (*handler.Handler)(nil), createShowHandler)
}

func Test_should_panic_if_no_port_was_found_on_create_show_handler(t *testing.T) {
	invalidPortMap := inbound.PortMap{
		inbound.PortInvalid: mockCreateShowService,
	}

	assert.Panics(t, func() {
		NewCreateShowHandler(invalidPortMap)
	})
}

func Test_should_return_route_on_create_show(t *testing.T) {
	var route = createShowHandler.GetRoute()

	var expectedRoute = &handler.Route{
		Method: "POST",
		Path:   "/show",
	}

	assert.Equal(t, expectedRoute, route)
}

func Test_should_call_service_on_create_show(t *testing.T) {
	defer mockCreateShowService.init()
	var createdShowDto *showResponseDto
	var context, recorder = handlerTestSetup.GetTestGinContext(t)

	test := struct {
		webCommand           string
		expectedPortCommand  *inbound.CreateShowCommand
		expectedPortResponse *inbound.CreateShowResponse
		expectedWebResponse  *showResponseDto
	}{
		`{"Title":"some title", "Slug":"some slug"}`,
		&inbound.CreateShowCommand{
			Title: "some title",
			Slug:  "some slug",
		},
		&inbound.CreateShowResponse{
			Id:    "some-id",
			Title: "Mocked Title",
			Slug:  "Mocked Slug",
		},
		&showResponseDto{
			Id:    "some-id",
			Title: "Mocked Title",
			Slug:  "Mocked Slug",
		},
	}

	mockCreateShowService.returnsOnCreateShow = test.expectedPortResponse

	context.Request = httptest.NewRequest("POST", "/show", bytes.NewBuffer([]byte(test.webCommand)))

	createShowHandler.Handle(context)

	var err = json.Unmarshal(recorder.Body.Bytes(), &createdShowDto)

	assert.Equal(t, 1, mockCreateShowService.called)
	assert.Equal(t, test.expectedPortCommand, mockCreateShowService.command)
	assert.Nil(t, err)
	assert.Empty(t, context.Errors)
	assert.Equal(t, test.expectedWebResponse, createdShowDto)
	assert.Equal(t, http.StatusCreated, recorder.Code)
}

func Test_should_propagate_error_on_create_show(t *testing.T) {
	defer mockCreateShowService.init()
	var context, _ = handlerTestSetup.GetTestGinContext(t)
	expectedError := errors.New("some error")

	test := struct {
		requestBody          string
		expectedPortResponse error
	}{
		`{"Title":"some title", "Slug":"some slug"}`,
		expectedError,
	}

	mockCreateShowService.failsWith = test.expectedPortResponse

	context.Request = httptest.NewRequest("POST", "/show", bytes.NewBuffer([]byte(test.requestBody)))

	createShowHandler.Handle(context)

	assert.NotEmpty(t, context.Errors)
	assert.Equal(t, expectedError, (*context.Errors[0]).Err)
}

func Test_abort_if_dto_is_invalid_on_create_show(t *testing.T) {
	defer mockCreateShowService.init()
	var context, recorder = handlerTestSetup.GetTestGinContext(t)

	test := struct {
		webCommand string
	}{
		`{"Bad":"dto"}`,
	}

	context.Request = httptest.NewRequest("POST", "/show", bytes.NewBuffer([]byte(test.webCommand)))

	createShowHandler.Handle(context)

	assert.NotEmpty(t, context.Errors)
	assert.Equal(t, 400, recorder.Code)
}
