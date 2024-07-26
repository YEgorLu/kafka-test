package messages

import (
	"os"

	"github.com/google/uuid"
	"github.com/segmentio/kafka-go"
)

const (
	Topic    = "messages"
	ErrTopic = "messages-error"
)

var (
	GroupID string
)

func init() {
	var err error
	GroupID, err = os.Hostname()
	if err != nil {
		panic(err)
	}
}

type MessagesMessage struct {
	kafka.Message
	Value MessagesValue
}

type MessagesValue struct {
	Id uuid.UUID
}

func (m MessagesMessage) Format() MessagesMessage {
	m.Message.Value = []byte(m.Value.Id.String())
	return m
}
