package error

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_show_already_exists_is_an_error(t *testing.T) {
	tests := map[string]struct {
		err                 error
		expectedErrorString string
	}{
		"ShowAlreadyExistsError": {
			NewShowAlreadyExistsError("some-name"),
			"show with title 'some-name' or given slug already exists",
		},

		"ShowNotFoundError": {
			NewShowNotFoundError("some-id"),
			"show with id 'some-id' does not exist",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			assert.Implements(t, (*error)(nil), test.err)
			assert.ErrorContains(t, test.err, test.expectedErrorString)
		})
	}

}
