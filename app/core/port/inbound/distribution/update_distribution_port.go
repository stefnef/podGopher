package distribution

type UpdateDistributionCommand struct {
	ShowId         string
	DistributionId string
	Title          *string
	Slug           *string
	Episodes       *[]string
}

type UpdateDistributionPort interface {
	UpdateDistribution(command *UpdateDistributionCommand) error
}
