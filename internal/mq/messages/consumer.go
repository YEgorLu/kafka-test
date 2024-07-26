package messages

import (
	"context"

	"github.com/YEgorLu/kafka-test/internal/config"
	"github.com/YEgorLu/kafka-test/internal/logger"
	"github.com/google/uuid"
	"github.com/segmentio/kafka-go"
)

type MessagesConsumer struct {
	*kafka.Reader
	log logger.Logger
}

func (c *MessagesConsumer) ReadMessage(ctx context.Context) (MessagesMessage, error) {
	msg, err := c.Reader.ReadMessage(ctx)
	if err != nil {
		c.log.Error("failed reading message ", err)
		return MessagesMessage{}, err
	}
	id, err := uuid.ParseBytes(msg.Value)
	if err != nil {
		c.log.Error("failed parsing uuid ", err, " ", msg.Value)
		return MessagesMessage{}, err
	}
	parsed := MessagesMessage{Message: msg, Value: MessagesValue{Id: id}}
	return parsed, err
}

func NewConsumer(log logger.Logger) *MessagesConsumer {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:     config.Mq.Brokers,
		Topic:       Topic,
		GroupID:     GroupID,
		Logger:      log,
		ErrorLogger: log,
	})
	return &MessagesConsumer{reader, log}
}
