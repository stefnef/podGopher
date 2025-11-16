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

type createEpisodeTestService struct {
	called                 int
	command                *inbound.CreateEpisodeCommand
	returnsOnCreateEpisode *inbound.CreateEpisodeResponse
	failsWith              error
}

func (s *createEpisodeTestService) init() {
	s.called = 0
	s.command = nil
	s.returnsOnCreateEpisode = nil
	s.failsWith = nil
}

func (s *createEpisodeTestService) CreateEpisode(command *inbound.CreateEpisodeCommand) (episode *inbound.CreateEpisodeResponse, err error) {
	s.called++
	s.command = command
	return s.returnsOnCreateEpisode, s.failsWith
}

var mockCreateEpisodeService = new(createEpisodeTestService)

var createEpisodeHandler = NewCreateEpisodeHandler(inbound.PortMap{
	inbound.CreateEpisode: mockCreateEpisodeService,
})

func Test_should_implement_handler_for_create_episode(t *testing.T) {
	assert.NotNil(t, createEpisodeHandler)
	assert.Implements(t, (*Handler)(nil), createEpisodeHandler)
}

func Test_should_panic_if_no_port_was_found_on_create_episode_handler(t *testing.T) {
	invalidPortMap := inbound.PortMap{
		inbound.PortInvalid: mockCreateEpisodeService,
	}

	assert.Panics(t, func() {
		NewCreateEpisodeHandler(invalidPortMap)
	})
}

func Test_should_return_route_on_create_episode(t *testing.T) {
	var route = createEpisodeHandler.GetRoute()

	var expectedRoute = &Route{
		Method: "POST",
		Path:   "/show/:showId/episode",
	}

	assert.Equal(t, expectedRoute, route)
}

func Test_should_propagate_error_on_create_episode(t *testing.T) {
	defer mockCreateEpisodeService.init()
	var context, _ = GetTestGinContext()
	expectedError := errors.New("some error")

	test := struct {
		paramShowId          string
		requestBody          string
		expectedPortResponse error
	}{
		`some-error-id`,
		`{"Title":"some title"}`,
		expectedError,
	}

	mockCreateEpisodeService.failsWith = test.expectedPortResponse

	context.Request = httptest.NewRequest("POST", "/show/"+test.paramShowId+"/episode", bytes.NewBuffer([]byte(test.requestBody)))

	createEpisodeHandler.Handle(context)

	assert.NotEmpty(t, context.Errors)
	assert.Equal(t, expectedError, (*context.Errors[0]).Err)
}

func Test_should_call_service_on_create_episode(t *testing.T) {
	defer mockCreateEpisodeService.init()
	var createEpisodeDto *episodeResponseDto
	var context, recorder = GetTestGinContext()

	test := struct {
		webParameterShowId   string
		webRequestBody       string
		expectedPortCommand  *inbound.CreateEpisodeCommand
		expectedPortResponse *inbound.CreateEpisodeResponse
		expectedWebResponse  *episodeResponseDto
	}{
		`some-show-id`,
		`{"Title":"some title"}`,
		&inbound.CreateEpisodeCommand{
			ShowId: "some-show-id",
			Title:  "some title",
		},
		&inbound.CreateEpisodeResponse{
			Id:     "some-id",
			ShowId: "Mocked Show Id",
			Title:  "Mocked Title",
		},
		&episodeResponseDto{
			Id:     "some-id",
			ShowId: "Mocked Show Id",
			Title:  "Mocked Title",
		},
	}

	mockCreateEpisodeService.returnsOnCreateEpisode = test.expectedPortResponse

	context.Request = httptest.NewRequest("POST", "/show/"+test.webParameterShowId+"/episode", bytes.NewBuffer([]byte(test.webRequestBody)))
	context.AddParam("showId", test.webParameterShowId)

	createEpisodeHandler.Handle(context)

	var err = json.Unmarshal(recorder.Body.Bytes(), &createEpisodeDto)

	assert.Equal(t, 1, mockCreateEpisodeService.called)
	assert.Equal(t, test.expectedPortCommand, mockCreateEpisodeService.command)
	assert.Nil(t, err)
	assert.Empty(t, context.Errors)
	assert.Equal(t, test.expectedWebResponse, createEpisodeDto)
	assert.Equal(t, http.StatusCreated, recorder.Code)
}
