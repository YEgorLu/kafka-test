package errors

import "github.com/google/uuid"

type NonExistsError struct {
	Id uuid.UUID
}

func (e NonExistsError) Error() string {
	return "Row with id " + e.Id.String() + " does not exists"
}
