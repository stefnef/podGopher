package repository

import "podGopher/core/port/outbound"

type InMemoryShowOutAdapter struct {
	shows []string
}

func (adapter *InMemoryShowOutAdapter) SaveShow(title string) (err error) {
	(*adapter).shows = append((*adapter).shows, title)
	return nil
}

func (adapter *InMemoryShowOutAdapter) ExistsByTitle(title string) bool {
	for _, value := range (*adapter).shows {
		if value == title {
			return true
		}
	}
	return false
}

func NewInMemoryShowRepository() outbound.SaveShowPort {
	return &InMemoryShowOutAdapter{shows: make([]string, 0)}
}
