package error

import "fmt"

type ShowAlreadyExistsError struct {
	Name string
}

type EpisodeAlreadyExistsError struct {
	Name string
}

type EpisodeNotFoundError struct {
	Id string
}

type ShowNotFoundError struct {
	Id string
}

type DistributionAlreadyExistsError struct {
	Name string
}

type DistributionNotFoundError struct {
	Id string
}

func (e ShowNotFoundError) Error() string {
	return fmt.Sprintf("show with id '%v' does not exist", e.Id)
}

func (e ShowAlreadyExistsError) Error() string {
	return fmt.Sprintf("show with title '%s' or given slug already exists", e.Name)
}

func (e EpisodeAlreadyExistsError) Error() string {
	return fmt.Sprintf("episode with title '%s' already exists", e.Name)
}

func (e EpisodeNotFoundError) Error() string {
	return fmt.Sprintf("episode with id '%v' does not exist", e.Id)
}

func (e DistributionAlreadyExistsError) Error() string {
	return fmt.Sprintf("distribution with title '%s' or given slug already exists", e.Name)
}

func (e DistributionNotFoundError) Error() string {
	return fmt.Sprintf("distribution with id '%v' does not exist", e.Id)
}

func NewShowAlreadyExistsError(name string) *ShowAlreadyExistsError {
	return &ShowAlreadyExistsError{name}
}

func NewShowNotFoundError(id string) *ShowNotFoundError {
	return &ShowNotFoundError{id}
}

func NewEpisodeAlreadyExistsError(name string) *EpisodeAlreadyExistsError {
	return &EpisodeAlreadyExistsError{name}
}

func NewEpisodeNotFoundError(id string) *EpisodeNotFoundError {
	return &EpisodeNotFoundError{id}
}

func NewDistributionAlreadyExistsError(name string) *DistributionAlreadyExistsError {
	return &DistributionAlreadyExistsError{name}
}

func NewDistributionNotFoundError(id string) *DistributionNotFoundError {
	return &DistributionNotFoundError{id}
}
