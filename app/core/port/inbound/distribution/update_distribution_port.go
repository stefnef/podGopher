package distribution

type UpdateDistributionCommand struct {
	ShowId string
	Title  *string
	Slug   *string
}

type UpdateDistributionResponse struct {
	Id     string
	ShowId string
	Title  string
	Slug   string
}

type UpdateDistributionPort interface {
	UpdateDistribution(command *UpdateDistributionCommand) (distribution *UpdateDistributionResponse, err error)
}
