package rss

import (
	"bytes"
	"encoding/xml"
	"errors"
	"net/http"
	"net/http/httptest"
	domainError "podGopher/core/domain/error"
	"podGopher/core/port/inbound"
	"podGopher/core/port/inbound/rss"
	"podGopher/integration/web/handler"
	"podGopher/integration/web/handler/handlerTestSetup"
	"testing"

	"github.com/stretchr/testify/assert"
)

type getRSSTestService struct {
	called          int
	command         *rss.GetRSSCommand
	returnsOnGetRSS *rss.GetRSSResponse
	failsWith       error
}

func (s *getRSSTestService) init() {
	s.called = 0
	s.command = nil
	s.returnsOnGetRSS = nil
	s.failsWith = nil
}

func (s *getRSSTestService) GetRSS(command *rss.GetRSSCommand) (*rss.GetRSSResponse, error) {
	s.called++
	s.command = command
	return s.returnsOnGetRSS, s.failsWith
}

var mockGetRSSService = new(getRSSTestService)

var getRSSHandler = NewGetRSSHandler(inbound.PortMap{
	inbound.GetRSS: mockGetRSSService,
})

func Test_should_implement_handler_for_get_rss(t *testing.T) {
	assert.NotNil(t, getRSSHandler)
	assert.Implements(t, (*handler.Handler)(nil), getRSSHandler)
}

func Test_should_panic_if_no_port_was_found_on_get_rss_handler(t *testing.T) {
	invalidPortMap := inbound.PortMap{
		inbound.PortInvalid: mockGetRSSService,
	}

	assert.Panics(t, func() {
		NewGetRSSHandler(invalidPortMap)
	})
}

func Test_should_return_route_on_get_rss(t *testing.T) {
	var route = getRSSHandler.GetRoute()

	var expectedRoute = &handler.Route{
		Method: "GET",
		Path:   "/rss/:showSlug/:distributionSlug",
	}

	assert.Equal(t, expectedRoute, route)
}

func Test_should_propagate_error_on_get_rss(t *testing.T) {
	defer mockGetRSSService.init()
	var context, _ = handlerTestSetup.GetTestGinContext(t)
	expectedError := errors.New("some error")

	test := struct {
		paramShowSlug         string
		paramDistributionSlug string
		expectedPortResponse  error
	}{
		`some-error-show-slug`,
		`some-error-distribution-slug`,
		expectedError,
	}

	mockGetRSSService.failsWith = test.expectedPortResponse

	context.Request = httptest.NewRequest("GET", "/rss/foo/bar", bytes.NewBuffer([]byte("")))
	context.AddParam("showSlug", test.paramShowSlug)
	context.AddParam("distributionSlug", test.paramDistributionSlug)

	getRSSHandler.Handle(context)

	assert.NotEmpty(t, context.Errors)
	assert.Equal(t, expectedError, (*context.Errors[0]).Err)
}

func Test_should_query_all_required_params_on_get_rss(t *testing.T) {
	tests := []struct {
		name       string
		paramKey   string
		paramValue string
	}{
		{
			name:       "missing show slug",
			paramKey:   "distributionSlug",
			paramValue: "dist-slug",
		},
		{
			name:       "missing distribution slug",
			paramKey:   "showSlug",
			paramValue: "show-slug",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			defer mockGetRSSService.init()
			var context, _ = handlerTestSetup.GetTestGinContext(t)

			context.Request = httptest.NewRequest("GET", "/rss///", bytes.NewBuffer([]byte("")))
			context.AddParam(test.paramKey, test.paramValue)

			getRSSHandler.Handle(context)

			assert.NotEmpty(t, context.Errors)
			assert.Equal(t, domainError.NewRSSFeedNotFoundError(test.paramValue), (*context.Errors[0]).Err)
		})
	}
}

func Test_should_call_service_on_get_rss(t *testing.T) {
	defer mockGetRSSService.init()
	var getRSSDto *rssResponseDto
	var context, recorder = handlerTestSetup.GetTestGinContext(t)

	test := struct {
		webParameterShowSlug         string
		webParameterDistributionSlug string
		expectedPortCommand          *rss.GetRSSCommand
		expectedPortResponse         *rss.GetRSSResponse
		expectedWebResponse          *rssResponseDto
	}{
		`show-slug`,
		`distribution-slug`,
		&rss.GetRSSCommand{
			ShowSlug:         "show-slug",
			DistributionSlug: "distribution-slug",
		},
		&rss.GetRSSResponse{
			ShowTitle:         "Mocked Shows's title",
			DistributionTitle: "Mocked distribution's title",
			Episodes: []*rss.GetRSSEpisodeResponse{
				{"Mocked Episode's title"},
			},
		},
		&rssResponseDto{
			ShowTitle:         "Mocked Shows's title",
			DistributionTitle: "Mocked distribution's title",
			Episodes: []rssEpisodeResponseDto{
				{
					Title: "Mocked Episode's title",
				},
			},
		},
	}

	mockGetRSSService.returnsOnGetRSS = test.expectedPortResponse

	context.Request = httptest.NewRequest("GET", "/rss/foo/bar", bytes.NewBuffer([]byte("")))
	context.AddParam("showSlug", test.webParameterShowSlug)
	context.AddParam("distributionSlug", test.webParameterDistributionSlug)

	getRSSHandler.Handle(context)

	if err := xml.Unmarshal(recorder.Body.Bytes(), &getRSSDto); err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, 1, mockGetRSSService.called)
	assert.Equal(t, test.expectedPortCommand, mockGetRSSService.command)
	assert.Empty(t, context.Errors)
	assert.Equal(t, test.expectedWebResponse, getRSSDto)
	assert.Equal(t, http.StatusOK, recorder.Code)
}
