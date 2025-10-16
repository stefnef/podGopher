package handler

import "podGopher/core/port/inbound"

type CreateShowHandler struct {
	route *Route
	port  inbound.CreateShowPort
}

type CreateShowCommand struct {
	Title string
}

func (h *CreateShowHandler) getRoute() *Route {
	return h.route
}

// TODO use responseWriter
func (h *CreateShowHandler) handle(command interface{}) {
	_ = h.port.CreateShow(&inbound.CreateShowCommand{Title: command.(*CreateShowCommand).Title})
}

func NewCreateShowHandler(portMap inbound.PortMap) *CreateShowHandler {
	return &CreateShowHandler{
		route: &Route{
			method: "POST",
			path:   "/show",
		},
		port: portMap[inbound.CreateShow].(inbound.CreateShowPort),
	}
}
