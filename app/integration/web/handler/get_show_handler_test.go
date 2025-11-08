package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"podGopher/core/port/inbound"
	"testing"

	"github.com/stretchr/testify/assert"
)

type getShowTestService struct {
	called           int
	command          *inbound.GetShowCommand
	returnsOnGetShow *inbound.GetShowResponse
	failsWith        error
}

func (s *getShowTestService) init() {
	s.called = 0
	s.command = nil
	s.returnsOnGetShow = nil
	s.failsWith = nil
}

func (s *getShowTestService) GetShow(command *inbound.GetShowCommand) (show *inbound.GetShowResponse, err error) {
	s.called++
	s.command = command
	return s.returnsOnGetShow, s.failsWith
}

var mockGetShowService = new(getShowTestService)

var getShowHandler = NewGetShowHandler(inbound.PortMap{
	inbound.GetShow: mockGetShowService,
})

func Test_should_implement_handler_for_get_show(t *testing.T) {
	assert.NotNil(t, getShowHandler)
	assert.Implements(t, (*Handler)(nil), getShowHandler)
}

func Test_should_panic_if_no_port_was_found_on_get_show_handler(t *testing.T) {
	invalidPortMap := inbound.PortMap{
		inbound.PortInvalid: mockCreateShowService,
	}

	assert.Panics(t, func() {
		NewGetShowHandler(invalidPortMap)
	})
}

func Test_should_return_route_on_get_show(t *testing.T) {
	var route = getShowHandler.GetRoute()

	var expectedRoute = &Route{
		Method: "GET",
		Path:   "/show/:showId",
	}

	assert.Equal(t, expectedRoute, route)
}

func Test_should_propagate_error_on_get_show(t *testing.T) {
	defer mockGetShowService.init()
	var context, _ = GetTestGinContext()
	expectedError := errors.New("some error")

	test := struct {
		paramShowId          string
		expectedPortResponse error
	}{
		`some-error-id`,
		expectedError,
	}

	mockGetShowService.failsWith = test.expectedPortResponse

	context.Request = httptest.NewRequest("GET", "/show/"+test.paramShowId, bytes.NewBuffer([]byte("")))

	getShowHandler.Handle(context)

	assert.NotEmpty(t, context.Errors)
	assert.Equal(t, expectedError, (*context.Errors[0]).Err)
}

func Test_should_call_service_on_get_show(t *testing.T) {
	defer mockGetShowService.init()
	var getShowDto *showResponseDto
	var context, recorder = GetTestGinContext()

	test := struct {
		webParameterShowId   string
		expectedPortCommand  *inbound.GetShowCommand
		expectedPortResponse *inbound.GetShowResponse
		expectedWebResponse  *showResponseDto
	}{
		`some-show-id`,
		&inbound.GetShowCommand{
			Id: "some-show-id",
		},
		&inbound.GetShowResponse{
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

	mockGetShowService.returnsOnGetShow = test.expectedPortResponse

	context.Request = httptest.NewRequest("GET", "/show/"+test.webParameterShowId, bytes.NewBuffer([]byte("")))
	context.AddParam("showId", test.webParameterShowId)

	getShowHandler.Handle(context)

	var err = json.Unmarshal(recorder.Body.Bytes(), &getShowDto)

	assert.Equal(t, 1, mockGetShowService.called)
	assert.Equal(t, test.expectedPortCommand, mockGetShowService.command)
	assert.Nil(t, err)
	assert.Empty(t, context.Errors)
	assert.Equal(t, test.expectedWebResponse, getShowDto)
	assert.Equal(t, http.StatusOK, recorder.Code)
}
