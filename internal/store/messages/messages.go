package messages

import (
	"context"
	"database/sql"

	"github.com/YEgorLu/kafka-test/internal/logger"
	"github.com/YEgorLu/kafka-test/internal/store/messages/models"
	"github.com/google/uuid"
)

type MessagesStore interface {
	Add(ctx context.Context, value string) (uuid.UUID, error)
	GetOne(ctx context.Context, id uuid.UUID) (models.Message, error)
	MarkProcessed(ctx context.Context, id uuid.UUID) error
	Delete(ctx context.Context, id uuid.UUID) error
	Stats(ctx context.Context) (models.Stats, error)
}

func NewStore(db *sql.DB, log logger.Logger) MessagesStore {
	return newPgMessagesStore(log, db)
}
