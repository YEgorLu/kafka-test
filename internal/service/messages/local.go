package messages

import (
	"context"
	"database/sql"
	"errors"

	"github.com/YEgorLu/kafka-test/internal/logger"
	mq "github.com/YEgorLu/kafka-test/internal/mq/messages"
	"github.com/YEgorLu/kafka-test/internal/service/messages/converters"
	serviceErrors "github.com/YEgorLu/kafka-test/internal/service/messages/errors"
	"github.com/YEgorLu/kafka-test/internal/service/messages/models"
	"github.com/YEgorLu/kafka-test/internal/store/messages"
	"github.com/google/uuid"
)

type localMessagesService struct {
	store       messages.MessagesStore
	producer    *mq.MessageProducer
	errProducer *mq.MessageProducer
	log         logger.Logger
}

func newLocalMessagesService(store messages.MessagesStore, mq *mq.MessageProducer, errMq *mq.MessageProducer, log logger.Logger) *localMessagesService {
	return &localMessagesService{
		store:       store,
		producer:    mq,
		errProducer: errMq,
		log:         log,
	}
}

func (s *localMessagesService) Add(ctx context.Context, msg models.NewMessage) (uuid.UUID, error) {
	id, err := s.store.Add(ctx, msg.Value)
	if err != nil {
		s.log.Error("can't add message with value ", msg.Value)
		return uuid.Nil, err
	}
	mqMsg := mq.MessagesMessage{Value: mq.MessagesValue{Id: id}}
	if err := s.producer.WriteMessages(ctx, mqMsg); err != nil {
		s.log.Error("can't place message ", id.String(), " in queue")
		go func() {
			err := s.store.Delete(context.Background(), id)
			if err != nil {
				s.log.Error("could not delete invalid message ", id.String())
			}
		}()
		return uuid.Nil, err
	}
	return id, nil
}

func (s *localMessagesService) Process(ctx context.Context, id uuid.UUID) error {
	var err error
	if err = s.store.MarkProcessed(ctx, id); err == nil {
		return nil
	}
	s.log.Error(err)
	if errors.Is(err, sql.ErrNoRows) {
		return serviceErrors.NonExistsError{Id: id}
	}
	if err := s.errProducer.WriteMessages(ctx, mq.MessagesMessage{Value: mq.MessagesValue{Id: id}}); err != nil {
		s.log.Error("error putting message in err queue ", id.String())
	}
	return err
}

func (s *localMessagesService) GetOne(ctx context.Context, id uuid.UUID) (models.Message, error) {
	storeMsg, err := s.store.GetOne(ctx, id)
	if err != nil {
		return models.Message{}, err
	}
	return converters.MsgToSv(storeMsg), nil
}

func (s *localMessagesService) Stats(ctx context.Context) (models.Stats, error) {
	storeStats, err := s.store.Stats(ctx)
	if err != nil {
		return models.Stats{}, err
	}
	return converters.StatsToSv(storeStats), nil
}
