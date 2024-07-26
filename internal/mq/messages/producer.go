package messages

import (
	"context"

	"github.com/YEgorLu/kafka-test/internal/config"
	"github.com/YEgorLu/kafka-test/internal/logger"
	kafka "github.com/segmentio/kafka-go"
)

type MessageProducer struct {
	*kafka.Writer
	log logger.Logger
}

func (c *MessageProducer) WriteMessages(ctx context.Context, data ...MessagesMessage) error {
	messages := make([]kafka.Message, 0, len(data))
	for _, v := range data {
		msg := v.Format()
		messages = append(messages, msg.Message)
	}
	return c.Writer.WriteMessages(ctx, messages...)
}

func NewMessagesProducer(log logger.Logger) *MessageProducer {
	writer := kafka.NewWriter(kafka.WriterConfig{
		Brokers:     config.Mq.Brokers,
		Topic:       Topic,
		Logger:      log,
		ErrorLogger: log,
	})
	return &MessageProducer{
		Writer: writer,
		log:    log,
	}
}

func NewErrorMessagesProducer(log logger.Logger) *MessageProducer {
	writer := kafka.NewWriter(kafka.WriterConfig{
		Brokers:      []string{"localhost:9092"},
		Topic:        ErrTopic,
		RequiredAcks: 1,
		Logger:       log,
		ErrorLogger:  log,
	})
	return &MessageProducer{
		Writer: writer,
		log:    log,
	}
}
