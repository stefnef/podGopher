package in_memory

import (
	"podGopher/core/port/outbound"

	"github.com/google/uuid"
)

type InMemoryShowOutAdapter struct {
	shows map[string]string
}

func (adapter *InMemoryShowOutAdapter) SaveShow(title string) (id string, err error) {
	id = uuid.NewString()
	adapter.shows[id] = title
	return id, nil
}

func (adapter *InMemoryShowOutAdapter) ExistsByTitle(title string) bool {
	for _, value := range adapter.shows {
		if value == title {
			return true
		}
	}
	return false
}

func NewInMemoryShowRepository() outbound.SaveShowPort {
	return &InMemoryShowOutAdapter{shows: make(map[string]string)}
}
