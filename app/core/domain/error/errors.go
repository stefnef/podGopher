package error

import "fmt"

type ShowAlreadyExistsError struct {
	Name string
}

type ShowNotFoundError struct {
	Id string
}

func (e ShowNotFoundError) Error() string {
	return fmt.Sprintf("show with id '%v' does not exist", e.Id)
}

func (e ShowAlreadyExistsError) Error() string {
	return fmt.Sprintf("show with title '%s' or given slug already exists", e.Name)
}

func NewShowAlreadyExistsError(name string) *ShowAlreadyExistsError {
	return &ShowAlreadyExistsError{name}
}

func NewShowNotFoundError(id string) *ShowNotFoundError {
	return &ShowNotFoundError{id}
}
