package handler

type Route struct {
	method string
	path   string
}

type Handler interface {
	getRoute() *Route
	handle(command interface{})
}

func CreateHandlers() []Handler { //TODO use portMap
	return []Handler{
		NewCreateShowHandler(nil),
	}
}
