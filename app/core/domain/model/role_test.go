package model

import (
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRoleNames(t *testing.T) {
	tests := []struct {
		role        ROLE
		expRoleName string
	}{
		{FOLLOWER,
			"FOLLOWER"},
		{
			EDITOR,
			"EDITOR",
		},
		{PRODUCER,
			"PRODUCER"},
	}
	for _, test := range tests {
		t.Run(test.expRoleName+" get Name", func(t *testing.T) {
			assert.Equal(t, test.expRoleName, test.role.Name())
		})

		t.Run(test.expRoleName+" to role", func(t *testing.T) {
			assert.Equal(t, test.role, ValueToRole(test.expRoleName))
		})
	}
}

func TestShouldSortShowRoleByRole(t *testing.T) {
	showRoles := []ShowRole{
		{ShowId: "editor", Role: EDITOR},
		{ShowId: "producer", Role: PRODUCER},
		{ShowId: "producer", Role: PRODUCER},
		{ShowId: "follower", Role: FOLLOWER},
	}

	expectedOrder := []ShowRole{
		{ShowId: "follower", Role: FOLLOWER},
		{ShowId: "producer", Role: PRODUCER},
		{ShowId: "producer", Role: PRODUCER},
		{ShowId: "editor", Role: EDITOR},
	}

	sort.Sort(ByRole(showRoles))

	assert.EqualValues(t, expectedOrder, showRoles)
}
