package show

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"podGopher/core/domain/model"
	"podGopher/core/port/inbound"
	"podGopher/core/port/inbound/show"
	"podGopher/integration/web/handler"
	"podGopher/integration/web/handler/handlerTestSetup"
	"testing"

	"github.com/stretchr/testify/assert"
)

var getAllShowsHandler = NewGetAllShowsHandler(inbound.PortMap{
	inbound.GetShow: mockGetShowService,
})

func Test_should_implement_handler_for_get_all_show(t *testing.T) {
	assert.NotNil(t, getAllShowsHandler)
	assert.Implements(t, (*handler.Handler)(nil), getAllShowsHandler)
}

func Test_should_panic_if_no_port_was_found_on_get_all_shows_handler(t *testing.T) {
	invalidPortMap := inbound.PortMap{
		inbound.PortInvalid: mockCreateShowService,
	}

	assert.Panics(t, func() {
		NewGetAllShowsHandler(invalidPortMap)
	})
}

func Test_should_return_route_on_get_all_shows(t *testing.T) {
	var route = getAllShowsHandler.GetRoute()

	var expectedRoute = &handler.Route{
		Method: "GET",
		Path:   "/show",
	}

	assert.Equal(t, expectedRoute, route)
}

func Test_should_propagate_error_on_get_all_shows(t *testing.T) {
	defer mockGetShowService.init()
	var context, _ = handlerTestSetup.GetTestGinContext(t)
	expectedError := errors.New("some error")

	mockGetShowService.failsWith = expectedError

	context.Request = httptest.NewRequest("GET", "/show", bytes.NewBuffer([]byte("")))

	getAllShowsHandler.Handle(context)

	assert.NotEmpty(t, context.Errors)
	assert.Equal(t, expectedError, (*context.Errors[0]).Err)
}

func Test_should_call_service_on_get_all_shows(t *testing.T) {
	defer mockGetShowService.init()
	var getAllShowsDto *allShowsResponseDto

	type testParameterStruct struct {
		title               string
		mockedPortResponse  *show.GetAllShowsResponse
		expectedWebResponse *allShowsResponseDto
	}

	tests := []testParameterStruct{
		{
			"one show returned",
			&show.GetAllShowsResponse{
				Shows: []*model.Show{
					{
						Id:       "some-id",
						Title:    "Mocked Title",
						Slug:     "Mocked Slug",
						Episodes: []string{},
					},
				},
			},
			&allShowsResponseDto{
				Shows: []allShowsItemResponseDto{
					{
						Id:    "some-id",
						Title: "Mocked Title",
					},
				},
			},
		},
		{
			"Nil response from port",
			&show.GetAllShowsResponse{
				Shows: nil,
			},
			&allShowsResponseDto{
				Shows: []allShowsItemResponseDto{},
			},
		},
		{
			"Empty response from port",
			&show.GetAllShowsResponse{
				Shows: []*model.Show{},
			},
			&allShowsResponseDto{
				Shows: []allShowsItemResponseDto{},
			},
		},
		{
			"2 shows returned",
			&show.GetAllShowsResponse{
				Shows: []*model.Show{
					{
						Id:       "some-id",
						Title:    "Mocked Title",
						Slug:     "Mocked Slug",
						Episodes: []string{},
					},
					{
						Id:       "2nd-show-id",
						Title:    "Some other Title",
						Slug:     "Some other Slug",
						Episodes: []string{"some-episode-id"},
					},
				},
			},
			&allShowsResponseDto{
				Shows: []allShowsItemResponseDto{
					{
						Id:    "some-id",
						Title: "Mocked Title",
					},
					{
						Id:    "2nd-show-id",
						Title: "Some other Title",
					},
				},
			},
		},
	}

	for _, tc := range tests {
		var context, recorder = handlerTestSetup.GetTestGinContext(t)
		mockGetShowService.init()

		t.Run(tc.title, func(t *testing.T) {
			mockGetShowService.returnsOnGetAllShow = tc.mockedPortResponse

			context.Request = httptest.NewRequest("GET", "/show", bytes.NewBuffer([]byte("")))

			getAllShowsHandler.Handle(context)

			var err = json.Unmarshal(recorder.Body.Bytes(), &getAllShowsDto)

			assert.Equal(t, 1, mockGetShowService.called)
			assert.Nil(t, err)
			assert.Empty(t, context.Errors)
			assert.Equal(t, tc.expectedWebResponse, getAllShowsDto)
			assert.Equal(t, http.StatusOK, recorder.Code)
		})
	}
}
