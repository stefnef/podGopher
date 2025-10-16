package inbound

type CreateShowCommand struct {
	Title string
}

type CreateShowResponse struct {
	Title string
}

type CreateShowPort interface {
	CreateShow(command *CreateShowCommand) (show *CreateShowResponse, err error)
}
