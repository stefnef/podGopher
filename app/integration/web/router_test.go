package web

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	domainError "podGopher/core/domain/error"
	"podGopher/core/domain/service/distribution"
	"podGopher/core/domain/service/episode"
	"podGopher/core/domain/service/show"
	"podGopher/core/port/inbound"
	inboundDistribution "podGopher/core/port/inbound/distribution"
	episode2 "podGopher/core/port/inbound/episode"
	show2 "podGopher/core/port/inbound/show"
	"testing"

	"github.com/stretchr/testify/assert"
)

type responseMock struct {
	Text      string `json:"Text"`
	failsWith error
}

var exampleRequests = map[string]string{
	"postShow":         `{"Title":"some title", "Slug":"some slug"}`,
	"postEpisode":      `{"Title":"some title"}`,
	"postDistribution": `{"Title":"some title", "Slug":"some slug"}`,
}

var response responseMock

type mockInboundPort struct{}

func (port *mockInboundPort) CreateShow(*show2.CreateShowCommand) (show *show2.CreateShowResponse, err error) {
	response.Text += "CreateShow"
	return &show2.CreateShowResponse{Title: "CreateShow"}, response.failsWith
}

func (port *mockInboundPort) GetShow(*show2.GetShowCommand) (show *show2.GetShowResponse, err error) {
	response.Text += "GetShow"
	return &show2.GetShowResponse{}, response.failsWith
}

func (port *mockInboundPort) CreateEpisode(*episode2.CreateEpisodeCommand) (episode *episode2.CreateEpisodeResponse, err error) {
	response.Text += "PostEpisode"
	return &episode2.CreateEpisodeResponse{}, response.failsWith
}

func (port *mockInboundPort) GetEpisode(*episode2.GetEpisodeCommand) (episode *episode2.GetEpisodeResponse, err error) {
	response.Text += "GetEpisode"
	return &episode2.GetEpisodeResponse{}, response.failsWith
}

func (port *mockInboundPort) CreateDistribution(*inboundDistribution.CreateDistributionCommand) (distribution *inboundDistribution.CreateDistributionResponse, err error) {
	response.Text += "PostDistribution"
	return &inboundDistribution.CreateDistributionResponse{}, response.failsWith
}

func (port *mockInboundPort) GetDistribution(*inboundDistribution.GetDistributionCommand) (distribution *inboundDistribution.GetDistributionResponse, err error) {
	response.Text += "GetDistribution"
	return &inboundDistribution.GetDistributionResponse{}, response.failsWith
}

var mockPort = new(mockInboundPort)
var router = NewRouter(inbound.PortMap{
	inbound.CreateShow:         mockPort,
	inbound.GetShow:            mockPort,
	inbound.CreateEpisode:      mockPort,
	inbound.GetEpisode:         mockPort,
	inbound.CreateDistribution: mockPort,
	inbound.GetDistribution:    mockPort,
})

func setup() {
	response = responseMock{Text: "", failsWith: nil}
}

func Test_should_return_NotFound_on_wrong_path(t *testing.T) {
	recorder := doRequest("GET", "/", "")

	assert.Equal(t, "404 page not found", recorder.Body.String())
	assert.Equal(t, http.StatusNotFound, recorder.Code)
}

func Test_should_post_a_show(t *testing.T) {
	setup()
	doRequest("POST", "/show", exampleRequests["postShow"])

	assert.Equal(t, "CreateShow", response.Text)
}

func Test_should_get_a_show(t *testing.T) {
	setup()
	doRequest("GET", "/show/some-show-id", "")

	assert.Equal(t, "GetShow", response.Text)
}

func Test_should_post_an_episode(t *testing.T) {
	setup()
	doRequest("POST", "/show/show-id/episode", exampleRequests["postEpisode"])

	assert.Equal(t, "PostEpisode", response.Text)
}

func Test_should_get_an_episode(t *testing.T) {
	setup()
	doRequest("GET", "/show/some-show-id/episode/some-episode-id", "")

	assert.Equal(t, "GetEpisode", response.Text)
}

func Test_should_post_a_distribution(t *testing.T) {
	setup()
	doRequest("POST", "/show/show-id/distribution", exampleRequests["postDistribution"])

	assert.Equal(t, "PostDistribution", response.Text)
}

func Test_should_get_a_distribution(t *testing.T) {
	setup()
	doRequest("GET", "/show/some-show-id/distribution/some-distribution-id", "")

	assert.Equal(t, "GetDistribution", response.Text)
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
		inbound.CreateEpisode:      episode.NewCreateEpisodeService(nil, nil),
		inbound.GetEpisode:         episode.NewGetEpisodeService(nil, nil),
		inbound.CreateDistribution: distribution.NewCreateDistributionService(nil, nil),
		inbound.GetDistribution:    distribution.NewGetDistributionService(nil, nil),
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
