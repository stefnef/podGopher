package model

type ShowRole struct {
	ShowId string
	Role   ROLE
}

type ROLE int

const (
	FOLLOWER ROLE = iota
	PRODUCER
	EDITOR
)
