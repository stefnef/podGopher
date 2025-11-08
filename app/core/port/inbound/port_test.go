package inbound

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_should_map_to_string(t *testing.T) {
	tests := map[string]struct {
		port           PORT
		expectedString string
	}{
		"INVALID": {
			PortInvalid,
			"INBOUND_PORT_INVALID",
		},
		"CREATE_SHOW": {
			CreateShow,
			"INBOUND_CREATE_SHOW",
		},
		"GET_SHOW": {
			GetShow,
			"INBOUND_GET_SHOW",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, test.expectedString, test.port.String())
		})
	}
}
