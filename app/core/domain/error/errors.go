package error

import "fmt"

type ShowAlreadyExistsError struct {
	Name string
}

func (e ShowAlreadyExistsError) Error() string {
	return fmt.Sprintf("show with title '%s' already exists", e.Name)
}

func NewShowAlreadyExistsError(name string) *ShowAlreadyExistsError {
	return &ShowAlreadyExistsError{name}
}
