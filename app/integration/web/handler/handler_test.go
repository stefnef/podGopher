package handler

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_should_create_handlers(t *testing.T) {
	var handlers = CreateHandlers()

	assert.NotEmpty(t, handlers)
	assert.Len(t, handlers, 1)
}
