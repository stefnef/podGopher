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
	assert.Implements(t, (*handler.Handler)(nil), getShowHandler)
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

	var expectedRoute = &handler.Route{
		Method: "GET",
		Path:   "/show/:showId",
	}

	assert.Equal(t, expectedRoute, route)
}

func Test_should_propagate_error_on_get_show(t *testing.T) {
	defer mockGetShowService.init()
	var context, _ = handlerTestSetup.GetTestGinContext(t)
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

	type testParameterStruct struct {
		webParameterShowId  string
		expectedPortCommand *inbound.GetShowCommand
		mockedPortResponse  *inbound.GetShowResponse
		expectedWebResponse *showResponseDto
	}

	tests := []testParameterStruct{
		{
			`some-show-without-episodes-id`,
			&inbound.GetShowCommand{
				Id: "some-show-without-episodes-id",
			},
			&inbound.GetShowResponse{
				Id:       "some-id",
				Title:    "Mocked Title",
				Slug:     "Mocked Slug",
				Episodes: []string{},
			},
			&showResponseDto{
				Id:       "some-id",
				Title:    "Mocked Title",
				Slug:     "Mocked Slug",
				Episodes: []string{},
			},
		},
		{
			`some-show-null-episodes`,
			&inbound.GetShowCommand{
				Id: "some-show-null-episodes",
			},
			&inbound.GetShowResponse{
				Id:       "some-id",
				Title:    "Mocked Title",
				Slug:     "Mocked Slug",
				Episodes: nil,
			},
			&showResponseDto{
				Id:       "some-id",
				Title:    "Mocked Title",
				Slug:     "Mocked Slug",
				Episodes: []string{},
			},
		},
		{
			`some-show-with-episodes-id`,
			&inbound.GetShowCommand{
				Id: "some-show-with-episodes-id",
			},
			&inbound.GetShowResponse{
				Id:       "some-id",
				Title:    "Mocked Title",
				Slug:     "Mocked Slug",
				Episodes: []string{"some-episode-id"},
			},
			&showResponseDto{
				Id:       "some-id",
				Title:    "Mocked Title",
				Slug:     "Mocked Slug",
				Episodes: []string{"some-episode-id"},
			},
		},
	}

	for _, tc := range tests {
		var context, recorder = handlerTestSetup.GetTestGinContext(t)
		mockGetShowService.init()

		t.Run(tc.webParameterShowId, func(t *testing.T) {
			mockGetShowService.returnsOnGetShow = tc.mockedPortResponse

			context.Request = httptest.NewRequest("GET", "/show/"+tc.webParameterShowId, bytes.NewBuffer([]byte("")))
			context.AddParam("showId", tc.webParameterShowId)

			getShowHandler.Handle(context)

			var err = json.Unmarshal(recorder.Body.Bytes(), &getShowDto)

			assert.Equal(t, 1, mockGetShowService.called)
			assert.Equal(t, tc.expectedPortCommand, mockGetShowService.command)
			assert.Nil(t, err)
			assert.Empty(t, context.Errors)
			assert.Equal(t, tc.expectedWebResponse, getShowDto)
			assert.Equal(t, http.StatusOK, recorder.Code)
		})
	}
}
