package handler

import (
	"podGopher/core/port/inbound"
	"testing"

	"github.com/stretchr/testify/assert"
)

type createShowTestService struct {
	called  int
	command *inbound.CreateShowCommand
}

func (s *createShowTestService) init() {
	s.called = 0
	s.command = nil
}

func (s *createShowTestService) CreateShow(command *inbound.CreateShowCommand) (err error) {
	s.called++
	s.command = command
	return nil
}

var mockCreateShowService = new(createShowTestService)
var createShowHandler = NewCreateShowHandler(mockCreateShowService)

func Test_should_implement_handler_for_create_show(t *testing.T) {
	assert.NotNil(t, createShowHandler)
	assert.Implements(t, (*Handler)(nil), createShowHandler)
}

func Test_should_return_route_on_create_show(t *testing.T) {
	var route = createShowHandler.getRoute()

	var expectedRoute = &Route{
		method: "POST",
		path:   "/show",
	}

	assert.Equal(t, expectedRoute, route)
}

func Test_should_call_service_on_create_show(t *testing.T) {
	defer mockCreateShowService.init()

	var webCommand = &CreateShowCommand{
		Title: "some title",
	}
	var expectedServiceCommand = &inbound.CreateShowCommand{
		Title: "some title",
	}

	createShowHandler.handle(webCommand)

	assert.Equal(t, mockCreateShowService.called, 1)
	assert.Equal(t, mockCreateShowService.command, expectedServiceCommand)
}
