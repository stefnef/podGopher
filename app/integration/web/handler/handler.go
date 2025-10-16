package handler

import "podGopher/core/port/inbound"

type Route struct {
	method string
	path   string
}

type Handler interface {
	getRoute() *Route
	handle(command interface{})
}

func CreateHandlers(portMap inbound.PortMap) []Handler {
	return []Handler{
		NewCreateShowHandler(portMap),
	}
}
