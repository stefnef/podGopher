package show

type GetShowCommand struct {
	Id string
}

type GetShowResponse struct {
	Id       string
	Title    string
	Slug     string
	Episodes []string
}

type GetShowPort interface {
	GetShow(command *GetShowCommand) (show *GetShowResponse, err error)
}
