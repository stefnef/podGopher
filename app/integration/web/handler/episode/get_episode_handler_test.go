package episode

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	error2 "podGopher/core/domain/error"
	"podGopher/core/port/inbound"
	"podGopher/integration/web/handler"
	"podGopher/integration/web/handler/handlerTestSetup"
	"testing"

	"github.com/stretchr/testify/assert"
)

type getEpisodeTestService struct {
	called              int
	command             *inbound.GetEpisodeCommand
	returnsOnGetEpisode *inbound.GetEpisodeResponse
	failsWith           error
}

func (s *getEpisodeTestService) init() {
	s.called = 0
	s.command = nil
	s.returnsOnGetEpisode = nil
	s.failsWith = nil
}

func (s *getEpisodeTestService) GetEpisode(command *inbound.GetEpisodeCommand) (episode *inbound.GetEpisodeResponse, err error) {
	s.called++
	s.command = command
	return s.returnsOnGetEpisode, s.failsWith
}

var mockGetEpisodeService = new(getEpisodeTestService)

var getEpisodeHandler = NewGetEpisodeHandler(inbound.PortMap{
	inbound.GetEpisode: mockGetEpisodeService,
})

func Test_should_implement_handler_for_get_episode(t *testing.T) {
	assert.NotNil(t, getEpisodeHandler)
	assert.Implements(t, (*handler.Handler)(nil), getEpisodeHandler)
}

func Test_should_panic_if_no_port_was_found_on_get_episode_handler(t *testing.T) {
	invalidPortMap := inbound.PortMap{
		inbound.PortInvalid: mockGetEpisodeService,
	}

	assert.Panics(t, func() {
		NewGetEpisodeHandler(invalidPortMap)
	})
}

func Test_should_return_route_on_get_episode(t *testing.T) {
	var route = getEpisodeHandler.GetRoute()

	var expectedRoute = &handler.Route{
		Method: "GET",
		Path:   "/show/:showId/episode/:episodeId",
	}

	assert.Equal(t, expectedRoute, route)
}

func Test_should_propagate_error_on_get_episode(t *testing.T) {
	defer mockGetEpisodeService.init()
	var context, _ = handlerTestSetup.GetTestGinContext(t)
	expectedError := errors.New("some error")

	test := struct {
		paramShowId          string
		paramEpisodeId       string
		expectedPortResponse error
	}{
		`some-error-show-id`,
		`some-error-episode-id`,
		expectedError,
	}

	mockGetEpisodeService.failsWith = test.expectedPortResponse

	context.Request = httptest.NewRequest("GET", "/show/"+test.paramShowId+"/episode/"+test.paramEpisodeId, bytes.NewBuffer([]byte("")))
	context.AddParam("showId", test.paramShowId)
	context.AddParam("episodeId", test.paramEpisodeId)

	getEpisodeHandler.Handle(context)

	assert.NotEmpty(t, context.Errors)
	assert.Equal(t, expectedError, (*context.Errors[0]).Err)
}

func Test_should_query_showId_param_on_get_episode(t *testing.T) {
	defer mockGetEpisodeService.init()
	var context, _ = handlerTestSetup.GetTestGinContext(t)

	context.Request = httptest.NewRequest("GET", "/show//episode/some-episode-id", bytes.NewBuffer([]byte("")))
	context.AddParam("episodeId", "some-episode-id")

	getEpisodeHandler.Handle(context)

	assert.NotEmpty(t, context.Errors)
	assert.Equal(t, error2.NewShowNotFoundError(""), (*context.Errors[0]).Err)
}

func Test_should_call_service_on_get_episode(t *testing.T) {
	defer mockGetEpisodeService.init()
	var getEpisodeDto *episodeResponseDto
	var context, recorder = handlerTestSetup.GetTestGinContext(t)

	test := struct {
		webParameterShowId    string
		webParameterEpisodeId string
		expectedPortCommand   *inbound.GetEpisodeCommand
		expectedPortResponse  *inbound.GetEpisodeResponse
		expectedWebResponse   *episodeResponseDto
	}{
		`some-show-id`,
		`some-episode-id`,
		&inbound.GetEpisodeCommand{
			ShowId:    "some-show-id",
			EpisodeId: "some-episode-id",
		},
		&inbound.GetEpisodeResponse{
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

	mockGetEpisodeService.returnsOnGetEpisode = test.expectedPortResponse

	context.Request = httptest.NewRequest("GET", "/show/"+test.webParameterShowId+"/episode/"+test.webParameterEpisodeId, bytes.NewBuffer([]byte("")))
	context.AddParam("showId", test.webParameterShowId)
	context.AddParam("episodeId", test.webParameterEpisodeId)

	getEpisodeHandler.Handle(context)

	var err = json.Unmarshal(recorder.Body.Bytes(), &getEpisodeDto)

	assert.Equal(t, 1, mockGetEpisodeService.called)
	assert.Equal(t, test.expectedPortCommand, mockGetEpisodeService.command)
	assert.Nil(t, err)
	assert.Empty(t, context.Errors)
	assert.Equal(t, test.expectedWebResponse, getEpisodeDto)
	assert.Equal(t, http.StatusOK, recorder.Code)
}
