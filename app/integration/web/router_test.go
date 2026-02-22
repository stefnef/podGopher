package web

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	domainError "podGopher/core/domain/error"
	"podGopher/core/domain/service/distribution"
	"podGopher/core/domain/service/episode"
	"podGopher/core/domain/service/rss"
	"podGopher/core/domain/service/show"
	"podGopher/core/port/inbound"
	inboundDistribution "podGopher/core/port/inbound/distribution"
	inboundEpisode "podGopher/core/port/inbound/episode"
	inboundRSS "podGopher/core/port/inbound/rss"
	inboundShow "podGopher/core/port/inbound/show"
	"testing"

	"github.com/stretchr/testify/assert"
)

type responseMock struct {
	Text      string `json:"Text"`
	failsWith error
}

var exampleRequests = map[string]string{
	"postShow":          `{"Title":"some title", "Slug":"some slug"}`,
	"postEpisode":       `{"Title":"some title"}`,
	"postDistribution":  `{"Title":"some title", "Slug":"some slug"}`,
	"patchDistribution": `{"Title":"some title", "Slug":"some slug", "Episodes":["episode-id"]}`,
}

var response responseMock

type mockInboundPort struct{}

func (port *mockInboundPort) CreateShow(*inboundShow.CreateShowCommand) (show *inboundShow.CreateShowResponse, err error) {
	response.Text += "CreateShow"
	return &inboundShow.CreateShowResponse{Title: "CreateShow"}, response.failsWith
}

func (port *mockInboundPort) GetShow(*inboundShow.GetShowCommand) (show *inboundShow.GetShowResponse, err error) {
	response.Text += "GetShow"
	return &inboundShow.GetShowResponse{}, response.failsWith
}

func (port *mockInboundPort) GetAllShows() (shows *inboundShow.GetAllShowsResponse, err error) {
	response.Text += "GetAllShows"
	return &inboundShow.GetAllShowsResponse{}, response.failsWith
}

func (port *mockInboundPort) CreateEpisode(*inboundEpisode.CreateEpisodeCommand) (episode *inboundEpisode.CreateEpisodeResponse, err error) {
	response.Text += "PostEpisode"
	return &inboundEpisode.CreateEpisodeResponse{}, response.failsWith
}

func (port *mockInboundPort) GetEpisode(*inboundEpisode.GetEpisodeCommand) (episode *inboundEpisode.GetEpisodeResponse, err error) {
	response.Text += "GetEpisode"
	return &inboundEpisode.GetEpisodeResponse{}, response.failsWith
}

func (port *mockInboundPort) CreateDistribution(*inboundDistribution.CreateDistributionCommand) (distribution *inboundDistribution.CreateDistributionResponse, err error) {
	response.Text += "PostDistribution"
	return &inboundDistribution.CreateDistributionResponse{}, response.failsWith
}

func (port *mockInboundPort) GetDistribution(*inboundDistribution.GetDistributionCommand) (distribution *inboundDistribution.GetDistributionResponse, err error) {
	response.Text += "GetDistribution"
	return &inboundDistribution.GetDistributionResponse{}, response.failsWith
}

func (port *mockInboundPort) UpdateDistribution(*inboundDistribution.UpdateDistributionCommand) error {
	response.Text += "UpdateDistribution"
	return response.failsWith
}

func (port *mockInboundPort) GetRSS(*inboundRSS.GetRSSCommand) (rssResponse *inboundRSS.GetRSSResponse, err error) {
	response.Text += "GetRSS"
	return &inboundRSS.GetRSSResponse{}, response.failsWith
}

var mockPort = new(mockInboundPort)
var router = NewRouter(inbound.PortMap{
	inbound.CreateShow:         mockPort,
	inbound.GetShow:            mockPort,
	inbound.CreateEpisode:      mockPort,
	inbound.GetEpisode:         mockPort,
	inbound.CreateDistribution: mockPort,
	inbound.GetDistribution:    mockPort,
	inbound.UpdateDistribution: mockPort,
	inbound.GetRSS:             mockPort,
	inbound.GetAllShows:        mockPort,
})

func setup() {
	response = responseMock{Text: "", failsWith: nil}
}

func Test_should_return_NotFound_on_wrong_path(t *testing.T) {
	recorder := doRequest("GET", "/", "")

	assert.Equal(t, "404 page not found", recorder.Body.String())
	assert.Equal(t, http.StatusNotFound, recorder.Code)
}

func Test_should_serve_routes(t *testing.T) {
	tests := map[string]struct {
		method                  string
		path                    string
		requestBody             string
		expectedMockHandlerCall string
	}{
		"Post show": {
			"POST",
			"/show",
			exampleRequests["postShow"],
			"CreateShow",
		},
		"Get single show": {
			"GET",
			"/show/some-show-id",
			"",
			"GetShow",
		},
		"Get all shows": {
			"GET",
			"/show",
			"",
			"GetAllShows",
		},
		"Post episode": {
			"POST",
			"/show/show-id/episode",
			exampleRequests["postEpisode"],
			"PostEpisode",
		},
		"Get episode": {
			"GET",
			"/show/show-id/episode/episode-id",
			"",
			"GetEpisode",
		},
		"Post distribution": {
			"POST",
			"/show/show-id/distribution",
			exampleRequests["postDistribution"],
			"PostDistribution",
		},
		"Get distribution": {
			"GET",
			"/show/show-id/distribution/some-distribution-id",
			"",
			"GetDistribution",
		},
		"Patch distribution": {
			"PATCH",
			"/show/show-id/distribution/some-distribution-id",
			exampleRequests["patchDistribution"],
			"UpdateDistribution",
		},
		"Get RSS": {
			"GET",
			"/rss/show-slug/distribution-slug",
			"",
			"GetRSS",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			setup()
			doRequest(test.method, test.path, test.requestBody)

			assert.Equal(t, test.expectedMockHandlerCall, response.Text)
		})
	}
}

func Test_should_handle_errors(t *testing.T) {
	setup()

	tests := map[string]struct {
		err          error
		expectedCode int
		expectedMsg  string
	}{
		"already_exists": {
			err:          &domainError.ModelError{Category: domainError.AlreadyExists, Context: "FAKE"},
			expectedCode: 400,
			expectedMsg:  "FAKE",
		},
		"not_found": {
			&domainError.ModelError{Category: domainError.NotFound, Context: "Not found FAKE"},
			404,
			"Not found FAKE",
		},
		"model_unknown": {
			&domainError.ModelError{Category: domainError.Unknown, Context: "Unknown FAKE"},
			500,
			"Unknown FAKE",
		},
		"data conflict": {&domainError.ModelError{Category: domainError.DataConflict, Context: "Some data conflict"},
			400,
			"Some data conflict",
		},
		"unknown": {
			errors.New("FAKE"),
			500,
			"Internal Server Error",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			response.failsWith = test.err

			recorder := doRequest("POST", "/show", exampleRequests["postShow"])

			assert.Equal(t, test.expectedCode, recorder.Code)
			assert.Contains(t, recorder.Body.String(), test.expectedMsg)
		})
	}
}

func Test_should_create_handlers(t *testing.T) {
	portMap := inbound.PortMap{
		inbound.CreateShow:         show.NewCreateShowService(nil),
		inbound.GetShow:            show.NewGetShowService(nil),
		inbound.GetAllShows:        show.NewGetShowService(nil),
		inbound.CreateEpisode:      episode.NewCreateEpisodeService(nil, nil),
		inbound.GetEpisode:         episode.NewGetEpisodeService(nil, nil),
		inbound.CreateDistribution: distribution.NewCreateDistributionService(nil, nil),
		inbound.GetDistribution:    distribution.NewGetDistributionService(nil, nil),
		inbound.UpdateDistribution: distribution.NewUpdateDistributionService(nil, nil, nil, nil),
		inbound.GetRSS:             rss.NewGetRSSService(nil),
	}

	var handlers = CreateHandlers(portMap)

	assert.NotEmpty(t, handlers)
	assert.Len(t, handlers, len(portMap))
}

func doRequest(method string, url string, requestBody string) *httptest.ResponseRecorder {
	recorder := httptest.NewRecorder()
	req, _ := http.NewRequest(method, url, bytes.NewBuffer([]byte(requestBody)))
	router.ServeHTTP(recorder, req)
	return recorder
}
