package distribution

type CreateDistributionCommand struct {
	ShowId string
	Title  string
	Slug   string
}

type CreateDistributionResponse struct {
	Id     string
	ShowId string
	Title  string
	Slug   string
}

type CreateDistributionPort interface {
	CreateDistribution(command *CreateDistributionCommand) (distribution *CreateDistributionResponse, err error)
}
