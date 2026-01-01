package distribution

import (
	"bytes"
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

type createDistributionTestService struct {
	called                      int
	command                     *distribution.CreateDistributionCommand
	returnsOnCreateDistribution *distribution.CreateDistributionResponse
	failsWith                   error
}

func (s *createDistributionTestService) init() {
	s.called = 0
	s.command = nil
	s.returnsOnCreateDistribution = nil
	s.failsWith = nil
}

func (s *createDistributionTestService) CreateDistribution(command *distribution.CreateDistributionCommand) (distribution *distribution.CreateDistributionResponse, err error) {
	s.called++
	s.command = command
	return s.returnsOnCreateDistribution, s.failsWith
}

var mockCreateDistributionService = new(createDistributionTestService)
var createDistributionHandler = NewCreateDistributionHandler(inbound.PortMap{
	inbound.CreateDistribution: mockCreateDistributionService,
})

func Test_should_implement_handler_for_create_distribution(t *testing.T) {
	assert.NotNil(t, createDistributionHandler)
	assert.Implements(t, (*handler.Handler)(nil), createDistributionHandler)
}

func Test_should_panic_if_no_port_was_found_on_create_distribution_handler(t *testing.T) {
	invalidPortMap := inbound.PortMap{
		inbound.PortInvalid: mockCreateDistributionService,
	}

	assert.Panics(t, func() {
		NewCreateDistributionHandler(invalidPortMap)
	})
}

func Test_should_return_route_on_create_distribution(t *testing.T) {
	var route = createDistributionHandler.GetRoute()

	var expectedRoute = &handler.Route{
		Method: "POST",
		Path:   "/show/:showId/distribution",
	}

	assert.Equal(t, expectedRoute, route)
}

func Test_should_call_service_on_create_distribution(t *testing.T) {
	defer mockCreateDistributionService.init()
	var createdDistributionDto *distributionResponseDto
	var context, recorder = handlerTestSetup.GetTestGinContext(t)

	test := struct {
		showId               string
		webCommand           string
		expectedPortCommand  *distribution.CreateDistributionCommand
		expectedPortResponse *distribution.CreateDistributionResponse
		expectedWebResponse  *distributionResponseDto
	}{
		"some-show-id",
		`{"title":"some title", "slug":"some slug"}`,
		&distribution.CreateDistributionCommand{
			ShowId: "some-show-id",
			Title:  "some title",
			Slug:   "some slug",
		},
		&distribution.CreateDistributionResponse{
			Id:     "some-id",
			ShowId: "some-show-id",
			Title:  "Mocked Title",
			Slug:   "Mocked Slug",
		},
		&distributionResponseDto{
			Id:       "some-id",
			ShowId:   "some-show-id",
			Title:    "Mocked Title",
			Slug:     "Mocked Slug",
			Episodes: []string{},
		},
	}

	mockCreateDistributionService.returnsOnCreateDistribution = test.expectedPortResponse

	context.Request = httptest.NewRequest("POST", "/show/"+test.showId+"/distribution", bytes.NewBuffer([]byte(test.webCommand)))
	context.AddParam("showId", test.showId)

	createDistributionHandler.Handle(context)

	var err = json.Unmarshal(recorder.Body.Bytes(), &createdDistributionDto)

	assert.Equal(t, 1, mockCreateDistributionService.called)
	assert.Equal(t, test.expectedPortCommand, mockCreateDistributionService.command)
	assert.Nil(t, err)
	assert.Empty(t, context.Errors)
	assert.Equal(t, test.expectedWebResponse, createdDistributionDto)
	assert.Equal(t, http.StatusCreated, recorder.Code)
}

func Test_should_propagate_error_on_create_distribution(t *testing.T) {
	defer mockCreateDistributionService.init()
	var context, _ = handlerTestSetup.GetTestGinContext(t)
	expectedError := errors.New("some error")

	test := struct {
		showId               string
		requestBody          string
		expectedPortResponse error
	}{
		"some-show-id",
		`{"title":"some title", "slug":"some slug"}`,
		expectedError,
	}

	mockCreateDistributionService.failsWith = test.expectedPortResponse

	context.Request = httptest.NewRequest("POST", "/show/"+test.showId+"/distribution", bytes.NewBuffer([]byte(test.requestBody)))
	context.AddParam("showId", test.showId)

	createDistributionHandler.Handle(context)

	assert.NotEmpty(t, context.Errors)
	assert.Equal(t, expectedError, (*context.Errors[0]).Err)
}

func Test_abort_if_dto_is_invalid_on_create_distribution(t *testing.T) {
	defer mockCreateDistributionService.init()
	var context, recorder = handlerTestSetup.GetTestGinContext(t)

	test := struct {
		showId     string
		webCommand string
	}{
		"some-show-id",
		`{"Bad":"dto"}`,
	}

	context.Request = httptest.NewRequest("POST", "/show/"+test.showId+"/distribution", bytes.NewBuffer([]byte(test.webCommand)))
	context.AddParam("showId", test.showId)

	createDistributionHandler.Handle(context)

	assert.NotEmpty(t, context.Errors)
	assert.Equal(t, 400, recorder.Code)
}
