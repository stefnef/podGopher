package inbound

type PORT int

type PortMap map[PORT]interface{}

const (
	PortInvalid PORT = iota
	CreateShow
	GetShow
	GetAllShows
	CreateEpisode
	GetEpisode
	CreateDistribution
	GetDistribution
	UpdateDistribution
	GetRSS
	CreateUser
	AssignUser //TODO use enum constant
)
