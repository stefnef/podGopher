package distribution

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	domainError "podGopher/core/domain/error"
	"podGopher/core/port/inbound"
	"podGopher/core/port/inbound/distribution"
	"podGopher/integration/web/handler"
	"podGopher/integration/web/handler/handlerTestSetup"
	"testing"

	"github.com/stretchr/testify/assert"
)

type updateDistributionTestService struct {
	called                      int
	command                     *distribution.UpdateDistributionCommand
	returnsOnUpdateDistribution *distribution.UpdateDistributionResponse
	failsWith                   error
}

func (s *updateDistributionTestService) init() {
	s.called = 0
	s.command = nil
	s.returnsOnUpdateDistribution = nil
	s.failsWith = nil
}

func (s *updateDistributionTestService) UpdateDistribution(command *distribution.UpdateDistributionCommand) (distribution *distribution.UpdateDistributionResponse, err error) {
	s.called++
	s.command = command
	return s.returnsOnUpdateDistribution, s.failsWith
}

var mockUpdateDistributionService = new(updateDistributionTestService)
var updateDistributionHandler = NewUpdateDistributionHandler(inbound.PortMap{
	inbound.UpdateDistribution: mockUpdateDistributionService,
})

func Test_should_implement_handler_for_update_distribution(t *testing.T) {
	assert.NotNil(t, updateDistributionHandler)
	assert.Implements(t, (*handler.Handler)(nil), updateDistributionHandler)
}

func Test_should_panic_if_no_port_was_found_on_update_distribution_handler(t *testing.T) {
	invalidPortMap := inbound.PortMap{
		inbound.PortInvalid: mockUpdateDistributionService,
	}

	assert.Panics(t, func() {
		NewUpdateDistributionHandler(invalidPortMap)
	})
}

func Test_should_return_route_on_update_distribution(t *testing.T) {
	var route = updateDistributionHandler.GetRoute()

	var expectedRoute = &handler.Route{
		Method: "PATCH",
		Path:   "/show/:showId/distribution/:distributionId",
	}

	assert.Equal(t, expectedRoute, route)
}

func ptrString(value string) *string {
	return &value
}

func Test_should_call_service_on_update_distribution(t *testing.T) {
	mockUpdateDistributionService.init()

	showId := "some-show-id"
	distributionId := "some-distribution-id"

	tests := map[string]struct {
		webCommand               string
		expectedPortCommandTitle *string
		expectedPortCommandSlug  *string
	}{
		"all fields are set": {
			`{"title":"some title", "slug":"some slug"}`,
			ptrString("some title"),
			ptrString("some slug"),
		},
		"title is set": {
			`{"title":"some title"}`,
			ptrString("some title"),
			nil,
		},
		"slug is set": {
			`{"slug":"some slug"}`,
			nil,
			ptrString("some slug"),
		},
	}

	for name, test := range tests {
		mockUpdateDistributionService.init()
		var context, recorder = handlerTestSetup.GetTestGinContext(t)

		t.Run(name, func(t *testing.T) {
			context.Request = httptest.NewRequest("PATCH", "/show/some-show-id"+"/distribution/some-distribution-id", bytes.NewBuffer([]byte(test.webCommand)))
			context.AddParam("showId", showId)
			context.AddParam("distributionId", distributionId)

			updateDistributionHandler.Handle(context)

			assert.Equal(t, 1, mockUpdateDistributionService.called)
			assert.Equal(t, test.expectedPortCommandTitle, mockUpdateDistributionService.command.Title)
			assert.Equal(t, test.expectedPortCommandSlug, mockUpdateDistributionService.command.Slug)
			assert.Nil(t, recorder.Body.Bytes())
			assert.Empty(t, context.Errors)
			assert.Equal(t, http.StatusNoContent, recorder.Code)
		})
	}
}

func Test_should_propagate_error_on_update_distribution(t *testing.T) {
	defer mockUpdateDistributionService.init()
	var context, _ = handlerTestSetup.GetTestGinContext(t)
	expectedError := errors.New("some error")

	showId := "some-show-id"
	distributionId := "some-distribtion-id"

	test := struct {
		requestBody          string
		expectedPortResponse error
	}{
		`{"title":"some title", "slug":"some slug"}`,
		expectedError,
	}

	mockUpdateDistributionService.failsWith = test.expectedPortResponse

	context.Request = httptest.NewRequest("PATCH", "/show/"+showId+"/distribution/"+distributionId, bytes.NewBuffer([]byte(test.requestBody)))
	context.AddParam("showId", showId)
	context.AddParam("distributionId", distributionId)

	updateDistributionHandler.Handle(context)

	assert.NotEmpty(t, context.Errors)
	assert.Equal(t, expectedError, (*context.Errors[0]).Err)
}

func Test_abort_if_dto_is_invalid_on_update_distribution(t *testing.T) {
	defer mockUpdateDistributionService.init()
	var context, recorder = handlerTestSetup.GetTestGinContext(t)

	test := struct {
		webCommand string
	}{
		`{"Bad":"dto"}`,
	}

	context.Request = httptest.NewRequest("PATCH", "/show/show-id/distribution/distribtion-id", bytes.NewBuffer([]byte(test.webCommand)))
	context.AddParam("showId", "show-id")
	context.AddParam("distributionId", "distribution-id")

	updateDistributionHandler.Handle(context)

	assert.NotEmpty(t, context.Errors)
	assert.Equal(t, "invalid request", context.Errors[0].Error())
	assert.Equal(t, 400, recorder.Code)
}

func Test_abort_if_showId_is_missing_on_update_distribution(t *testing.T) {
	defer mockUpdateDistributionService.init()

	var context, _ = handlerTestSetup.GetTestGinContext(t)

	context.Request = httptest.NewRequest("PATCH", "/show/show-id/distribution/distribtion-id", bytes.NewBuffer([]byte((`{"title":"some"}`))))
	context.AddParam("distributionId", "some-value")

	updateDistributionHandler.Handle(context)

	assert.NotEmpty(t, context.Errors)
	assert.Equal(t, domainError.NewShowNotFoundError(""), (*context.Errors[0]).Err)

}
