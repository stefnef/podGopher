package inbound

type CreateShowCommand struct {
	Title string
}

type CreateShowPort interface {
	CreateShow(command *CreateShowCommand) (err error)
}
