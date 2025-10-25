package error

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_show_already_exists_is_an_error(t *testing.T) {
	err := NewShowAlreadyExistsError("some-name")

	assert.Implements(t, (*error)(nil), err)
	assert.ErrorContains(t, err, "show with title 'some-name' already exists")
}
