package messages

import (
	"context"

	"github.com/YEgorLu/kafka-test/internal/logger"
	mq "github.com/YEgorLu/kafka-test/internal/mq/messages"
	"github.com/YEgorLu/kafka-test/internal/service/messages/models"
	"github.com/YEgorLu/kafka-test/internal/store/messages"
	"github.com/google/uuid"
)

type MessageService interface {
	Add(ctx context.Context, msg models.NewMessage) (uuid.UUID, error)
	Process(ctx context.Context, id uuid.UUID) error
	GetOne(ctx context.Context, id uuid.UUID) (models.Message, error)
	Stats(ctx context.Context) (models.Stats, error)
}

func NewMessagesService(store messages.MessagesStore, mq *mq.MessageProducer, errMq *mq.MessageProducer, log logger.Logger) *localMessagesService {
	return newLocalMessagesService(store, mq, errMq, log)
}
