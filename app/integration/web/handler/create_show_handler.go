package handler

import "podGopher/core/port/inbound"

type CreateShowHandler struct {
	route   *Route
	service inbound.CreateShowPort
}

type CreateShowCommand struct {
	Title string
}

func (h *CreateShowHandler) getRoute() *Route {
	return h.route
}

// TODO use responseWriter
func (h *CreateShowHandler) handle(command interface{}) {
	_ = h.service.CreateShow(&inbound.CreateShowCommand{Title: command.(*CreateShowCommand).Title})
}

func NewCreateShowHandler(service inbound.CreateShowPort) *CreateShowHandler {
	return &CreateShowHandler{
		route: &Route{
			method: "POST",
			path:   "/show",
		},
		service: service,
	}
}
