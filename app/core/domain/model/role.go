package model

type ShowRole struct {
	ShowId string
	Role   ROLE
}
type ByRole []ShowRole

func (b ByRole) Len() int {
	return len(b)
}

func (b ByRole) Less(i, j int) bool {
	return b[i].Role < b[j].Role
}

func (b ByRole) Swap(i, j int) {
	b[i], b[j] = b[j], b[i]
}

type ROLE int

func ValueToRole(name string) ROLE {
	return valueToRole[name]
}

var roleNames = map[ROLE]string{
	FOLLOWER: "FOLLOWER",
	PRODUCER: "PRODUCER",
	EDITOR:   "EDITOR",
}

var valueToRole = map[string]ROLE{
	"FOLLOWER": FOLLOWER,
	"PRODUCER": PRODUCER,
	"EDITOR":   EDITOR,
}

func (r ROLE) Name() string {
	return roleNames[r]
}

const (
	FOLLOWER ROLE = iota
	PRODUCER
	EDITOR
)
