package domainRole

import (
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_should_map_roles(t *testing.T) {
	tests := []struct {
		role        ROLE
		expRoleName string
	}{
		{ADMIN,
			"ADMIN",
		},
		{FOLLOWER,
			"FOLLOWER",
		},
		{
			EDITOR,
			"EDITOR",
		},
		{PRODUCER,
			"PRODUCER",
		},
		{UNKNOWN,
			"UNKNOWN",
		},
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

func Test_should_ignore_case_on_value_to_role(t *testing.T) {
	assert.Equal(t, FOLLOWER, ValueToRole("FolLOwer"))
}

func Test_should_map_to_unknown(t *testing.T) {
	assert.Equal(t, UNKNOWN, ValueToRole("bad role name"))
}

func Test_should_sort_show_role_by_role(t *testing.T) {
	showRoles := []ShowRole{
		{ShowId: "admin", Role: ADMIN},
		{ShowId: "editor", Role: EDITOR},
		{ShowId: "producer", Role: PRODUCER},
		{ShowId: "producer", Role: PRODUCER},
		{ShowId: "follower", Role: FOLLOWER},
		{ShowId: "FOLLOWER", Role: FOLLOWER},
		{ShowId: "unknown", Role: UNKNOWN},
	}

	expectedOrder := []ShowRole{
		{ShowId: "unknown", Role: UNKNOWN},
		{ShowId: "follower", Role: FOLLOWER},
		{ShowId: "FOLLOWER", Role: FOLLOWER},
		{ShowId: "editor", Role: EDITOR},
		{ShowId: "producer", Role: PRODUCER},
		{ShowId: "producer", Role: PRODUCER},
		{ShowId: "admin", Role: ADMIN},
	}

	sort.Sort(ByRole(showRoles))

	assert.EqualValues(t, expectedOrder, showRoles)
}
