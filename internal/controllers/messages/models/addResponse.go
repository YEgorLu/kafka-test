package models

import "github.com/google/uuid"

type AddResponse struct {
	Id uuid.UUID `json:"id"`
}
