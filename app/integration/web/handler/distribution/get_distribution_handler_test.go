package distribution

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"podGopher/core/port/inbound"
	"podGopher/core/port/inbound/distribution"
	"podGopher/integration/web/handler"
	"podGopher/integration/web/handler/handlerTestSetup"
	"testing"

	"github.com/stretchr/testify/assert"
)

type getDistributionTestService struct {
	called                   int
	command                  *distribution.GetDistributionCommand
	returnsOnGetDistribution *distribution.GetDistributionResponse
	failsWith                error
}

func (s *getDistributionTestService) init() {
	s.called = 0
	s.command = nil
	s.returnsOnGetDistribution = nil
	s.failsWith = nil
}

func (s *getDistributionTestService) GetDistribution(command *distribution.GetDistributionCommand) (distribution *distribution.GetDistributionResponse, err error) {
	s.called++
	s.command = command
	return s.returnsOnGetDistribution, s.failsWith
}

var mockGetDistributionService = new(getDistributionTestService)
var getDistributionHandler = NewGetDistributionHandler(inbound.PortMap{
	inbound.GetDistribution: mockGetDistributionService,
})

func Test_should_implement_handler_for_get_distribution(t *testing.T) {
	assert.NotNil(t, getDistributionHandler)
	assert.Implements(t, (*handler.Handler)(nil), getDistributionHandler)
}

func Test_should_panic_if_no_port_was_found_on_get_distribution_handler(t *testing.T) {
	invalidPortMap := inbound.PortMap{
		inbound.PortInvalid: mockGetDistributionService,
	}

	assert.Panics(t, func() {
		NewGetDistributionHandler(invalidPortMap)
	})
}

func Test_should_return_route_on_get_distribution(t *testing.T) {
	var route = getDistributionHandler.GetRoute()

	var expectedRoute = &handler.Route{
		Method: "GET",
		Path:   "/show/:showId/distribution/:distributionId",
	}

	assert.Equal(t, expectedRoute, route)
}

func Test_should_call_service_on_get_distribution(t *testing.T) {
	defer mockGetDistributionService.init()
	var foundDistributionDto *distributionResponseDto

	type testParamStruct struct {
		showId               string
		distributionId       string
		expectedPortCommand  *distribution.GetDistributionCommand
		expectedPortResponse *distribution.GetDistributionResponse
		expectedWebResponse  *distributionResponseDto
	}

	tests := []testParamStruct{
		{
			"some-show-id",
			"some-distribution-id",
			&distribution.GetDistributionCommand{
				ShowId:         "some-show-id",
				DistributionId: "some-distribution-id",
			},
			&distribution.GetDistributionResponse{
				Id:       "some-distribution-id",
				ShowId:   "some-show-id",
				Title:    "Mocked Title",
				Slug:     "Mocked Slug",
				Episodes: []string{"episode-1", "episode-2"},
			},
			&distributionResponseDto{
				Id:       "some-distribution-id",
				ShowId:   "some-show-id",
				Title:    "Mocked Title",
				Slug:     "Mocked Slug",
				Episodes: []string{"episode-1", "episode-2"},
			},
		},
		{
			"some-show-id",
			"some-distribution-id-with-empty-episodes",
			&distribution.GetDistributionCommand{
				ShowId:         "some-show-id",
				DistributionId: "some-distribution-id-with-empty-episodes",
			},
			&distribution.GetDistributionResponse{
				Id:       "some-distribution-id-with-empty-episodes",
				ShowId:   "some-show-id",
				Title:    "Mocked Title",
				Slug:     "Mocked Slug",
				Episodes: []string{},
			},
			&distributionResponseDto{
				Id:       "some-distribution-id-with-empty-episodes",
				ShowId:   "some-show-id",
				Title:    "Mocked Title",
				Slug:     "Mocked Slug",
				Episodes: []string{},
			},
		},
		{
			"some-show-id",
			"some-distribution-id-with-nil-episodes",
			&distribution.GetDistributionCommand{
				ShowId:         "some-show-id",
				DistributionId: "some-distribution-id-with-nil-episodes",
			},
			&distribution.GetDistributionResponse{
				Id:       "some-distribution-id-with-nil-episodes",
				ShowId:   "some-show-id",
				Title:    "Mocked Title",
				Slug:     "Mocked Slug",
				Episodes: nil,
			},
			&distributionResponseDto{
				Id:       "some-distribution-id-with-nil-episodes",
				ShowId:   "some-show-id",
				Title:    "Mocked Title",
				Slug:     "Mocked Slug",
				Episodes: []string{},
			},
		},
	}

	for _, test := range tests {
		var context, recorder = handlerTestSetup.GetTestGinContext(t)
		mockGetDistributionService.init()

		t.Run(test.distributionId, func(t *testing.T) {
			mockGetDistributionService.returnsOnGetDistribution = test.expectedPortResponse

			context.Request = httptest.NewRequest("GET", "/show/"+test.showId+"/distribution/"+test.distributionId, nil)
			context.AddParam("showId", test.showId)
			context.AddParam("distributionId", test.distributionId)

			getDistributionHandler.Handle(context)

			var err = json.Unmarshal(recorder.Body.Bytes(), &foundDistributionDto)

			assert.Equal(t, 1, mockGetDistributionService.called)
			assert.Equal(t, test.expectedPortCommand, mockGetDistributionService.command)
			assert.Nil(t, err)
			assert.Empty(t, context.Errors)
			assert.Equal(t, test.expectedWebResponse, foundDistributionDto)
			assert.Equal(t, http.StatusOK, recorder.Code)
		})
	}
}

func Test_should_propagate_error_on_get_distribution(t *testing.T) {
	defer mockGetDistributionService.init()
	var context, _ = handlerTestSetup.GetTestGinContext(t)
	expectedError := errors.New("some error")

	test := struct {
		showId               string
		distributionId       string
		expectedPortResponse error
	}{
		"some-show-id",
		"some-distribution-id",
		expectedError,
	}

	mockGetDistributionService.failsWith = test.expectedPortResponse

	context.Request = httptest.NewRequest("GET", "/show/"+test.showId+"/distribution/"+test.distributionId, nil)
	context.AddParam("showId", test.showId)
	context.AddParam("distributionId", test.distributionId)

	getDistributionHandler.Handle(context)

	assert.NotEmpty(t, context.Errors)
	assert.Equal(t, expectedError, (*context.Errors[0]).Err)
}
