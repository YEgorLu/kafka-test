package models

import (
	"time"

	"github.com/google/uuid"
)

type Message struct {
	Id          uuid.UUID
	Value       string
	Processed   bool
	DateCreate  time.Time
	DateProcess *time.Time
}

type NewMessage struct {
	Value string
}
