package web

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	error2 "podGopher/core/domain/error"
	"podGopher/core/domain/service/distribution"
	"podGopher/core/domain/service/episode"
	"podGopher/core/domain/service/show"
	"podGopher/core/port/inbound"
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

func (port *mockInboundPort) CreateShow(*inbound.CreateShowCommand) (show *inbound.CreateShowResponse, err error) {
	response.Text += "CreateShow"
	return &inbound.CreateShowResponse{Title: "CreateShow"}, response.failsWith
}

func (port *mockInboundPort) GetShow(*inbound.GetShowCommand) (show *inbound.GetShowResponse, err error) {
	response.Text += "GetShow"
	return &inbound.GetShowResponse{}, response.failsWith
}

func (port *mockInboundPort) CreateEpisode(*inbound.CreateEpisodeCommand) (episode *inbound.CreateEpisodeResponse, err error) {
	response.Text += "PostEpisode"
	return &inbound.CreateEpisodeResponse{}, response.failsWith
}

func (port *mockInboundPort) GetEpisode(*inbound.GetEpisodeCommand) (episode *inbound.GetEpisodeResponse, err error) {
	response.Text += "GetEpisode"
	return &inbound.GetEpisodeResponse{}, response.failsWith
}

func (port *mockInboundPort) CreateDistribution(*inbound.CreateDistributionCommand) (distribution *inbound.CreateDistributionResponse, err error) {
	response.Text += "PostDistribution"
	return &inbound.CreateDistributionResponse{}, response.failsWith
}

func (port *mockInboundPort) GetDistribution(*inbound.GetDistributionCommand) (distribution *inbound.GetDistributionResponse, err error) {
	response.Text += "GetDistribution"
	return &inbound.GetDistributionResponse{}, response.failsWith
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
			err:          &error2.ModelError{Category: error2.AlreadyExists, Context: "FAKE"},
			expectedCode: 400,
			expectedMsg:  "FAKE",
		},
		"not_found": {
			&error2.ModelError{Category: error2.NotFound, Context: "Not found FAKE"},
			404,
			"Not found FAKE",
		},
		"model_unknown": {
			&error2.ModelError{Category: error2.Unknown, Context: "Unknown FAKE"},
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
