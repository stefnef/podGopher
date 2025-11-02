package web

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	error2 "podGopher/core/domain/error"
	"podGopher/core/port/inbound"
	"testing"

	"github.com/stretchr/testify/assert"
)

type responseMock struct {
	Text      string `json:"Text"`
	failsWith error
}

var response responseMock

type mockInboundPort struct{}

func (port *mockInboundPort) CreateShow(*inbound.CreateShowCommand) (show *inbound.CreateShowResponse, err error) {
	response.Text += "CreateShow"
	return &inbound.CreateShowResponse{Title: "CreateShow"}, response.failsWith
}

var mockPort = new(mockInboundPort)
var router = NewRouter(inbound.PortMap{
	inbound.CreateShow: mockPort,
})

func setup() {
	response = responseMock{Text: "", failsWith: nil}
}

func Test_should_return_NotFound_on_wrong_path(t *testing.T) {
	recorder := doRequest("GET", "/")

	assert.Equal(t, "404 page not found", recorder.Body.String())
	assert.Equal(t, http.StatusNotFound, recorder.Code)
}

func Test_should_post_a_show(t *testing.T) {
	setup()
	doRequest("POST", "/show")

	assert.Equal(t, "CreateShow", response.Text)
}

func Test_should_handle_errors(t *testing.T) {
	setup()

	tests := map[string]struct {
		err          error
		expectedCode int
		expectedMsg  string
	}{
		"show_already_exists": {
			error2.NewShowAlreadyExistsError("FAKE"),
			400,
			"FAKE",
		},
		"unknown": {
			errors.New("FAKE"),
			500,
			"FAKE",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			response.failsWith = test.err

			recorder := doRequest("POST", "/show")

			assert.Equal(t, test.expectedCode, recorder.Code)
			assert.Contains(t, recorder.Body.String(), test.expectedMsg)
		})
	}
}

func doRequest(method string, url string) *httptest.ResponseRecorder {
	recorder := httptest.NewRecorder()
	postShowRequest := `{"Title":"some title", "Slug":"some slug"}`
	req, _ := http.NewRequest(method, url, bytes.NewBuffer([]byte(postShowRequest)))
	router.ServeHTTP(recorder, req)
	return recorder
}
