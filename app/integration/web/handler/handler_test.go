package handler

import (
	"podGopher/core/domain/service"
	"podGopher/core/port/inbound"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_should_create_handlers(t *testing.T) {
	portMap := inbound.PortMap{
		inbound.CreateShow: service.NewCreateShowService(nil),
	}

	var handlers = CreateHandlers(portMap)

	assert.NotEmpty(t, handlers)
	assert.Len(t, handlers, 1)
}
