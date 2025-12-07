package inbound

type GetDistributionCommand struct {
	DistributionId string
	ShowId         string
}

type GetDistributionResponse struct {
	Id       string
	ShowId   string
	Title    string
	Slug     string
	Episodes []string
}

type GetDistributionPort interface {
	GetDistribution(command *GetDistributionCommand) (distribution *GetDistributionResponse, err error)
}
