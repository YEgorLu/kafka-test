package models

import (
	"time"

	"github.com/google/uuid"
)

type GetOneResponse struct {
	Id          uuid.UUID  `json:"id"`
	Value       string     `json:"value"`
	Processed   bool       `json:"processed"`
	DateCreate  time.Time  `json:"dateCreate"`
	DateProcess *time.Time `json:"dateProcess,omitempty"`
}
