package inbound

type PORT int

type PortMap map[PORT]interface{}

const (
	PortInvalid PORT = iota
	CreateShow
	GetShow
	CreateEpisode
)
