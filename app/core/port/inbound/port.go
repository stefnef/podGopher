package inbound

type PORT int

type PortMap map[PORT]interface{}

const (
	PortInvalid PORT = iota
	CreateShow
	GetShow
)

func (p PORT) String() string {
	switch p {
	case CreateShow:
		return "INBOUND_CREATE_SHOW"
	case GetShow:
		return "INBOUND_GET_SHOW"
	default:
		return "INBOUND_PORT_INVALID"
	}
}
