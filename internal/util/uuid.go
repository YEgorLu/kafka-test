package util

import (
	"errors"

	"github.com/google/uuid"
)

func ParseUUID(idStr string) (uuid.UUID, error) {
	if idStr == "" {
		return uuid.Nil, errors.New("empty id")
	}
	id, err := uuid.Parse(idStr)
	if err != nil {
		return uuid.Nil, err
	}
	return id, nil
}
