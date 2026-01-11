package show

type CreateShowCommand struct {
	Title string
	Slug  string
}

type CreateShowResponse struct {
	Id    string
	Title string
	Slug  string
}

type CreateShowPort interface {
	CreateShow(command *CreateShowCommand) (show *CreateShowResponse, err error)
}
