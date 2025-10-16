package inbound

type PORT int

type PortMap map[PORT]interface{}

const (
	PortInvalid PORT = iota
	CreateShow
)

func (p PORT) String() string {
	switch p {
	case CreateShow:
		return "INBOUND_CREATE_SHOW"
	default:
		return "INBOUND_PORT_INVALID"
	}
}
