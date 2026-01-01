package error

import "fmt"

type Category int

const (
	Unknown Category = iota
	AlreadyExists
	NotFound
)

type ModelError struct {
	ModelCategory Category
	Context       string
}

func (e ModelError) Error() string {
	return e.Context
}

func (e ModelError) Category() Category {
	return e.ModelCategory
}

func NewShowAlreadyExistsError(name string) *ModelError {
	context := fmt.Sprintf("show with title '%v' or given slug already exists", name)
	return &ModelError{AlreadyExists, context}
}

func NewShowNotFoundError(id string) *ModelError {
	context := fmt.Sprintf("show with id '%v' does not exist", id)
	return &ModelError{NotFound, context}
}

func NewEpisodeAlreadyExistsError(name string) *ModelError {
	context := fmt.Sprintf("episode with title '%v' already exists", name)
	return &ModelError{AlreadyExists, context}
}

func NewEpisodeNotFoundError(id string) *ModelError {
	context := fmt.Sprintf("episode with id '%v' does not exist", id)
	return &ModelError{NotFound, context}
}

func NewDistributionAlreadyExistsError(name string) *ModelError {
	context := fmt.Sprintf("distribution with title '%v' or given slug already exists", name)
	return &ModelError{AlreadyExists, context}
}

func NewDistributionNotFoundError(id string) *ModelError {
	context := fmt.Sprintf("distribution with id '%v' does not exist", id)
	return &ModelError{NotFound, context}
}
