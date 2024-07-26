package messages

import (
	"context"
	"database/sql"
	"sync"
	"sync/atomic"
	"time"

	"golang.org/x/sync/errgroup"

	"github.com/YEgorLu/kafka-test/internal/logger"
	"github.com/YEgorLu/kafka-test/internal/store/messages/models"
	"github.com/google/uuid"
)

type pgMessagesStore struct {
	db             *sql.DB
	log            logger.Logger
	countAll       atomic.Int64
	countProcessed atomic.Int64
	loadCountMu    sync.Mutex
	countsLoaded   bool
}

func newPgMessagesStore(log logger.Logger, db *sql.DB) *pgMessagesStore {
	store := &pgMessagesStore{
		db:  db,
		log: log,
	}

	go func() {
		for range time.NewTicker(time.Hour).C {
			store.loadCountMu.Lock()
			store.countsLoaded = false
			store.loadCountMu.Unlock()
		}
	}()

	return store
}

func (s *pgMessagesStore) Add(ctx context.Context, value string) (uuid.UUID, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		s.log.Error("Error adding message", err)
		return uuid.Nil, err
	}
	defer tx.Rollback()
	row := tx.QueryRowContext(ctx, "INSERT INTO Messages(value, date_create) VALUES ($1, $2) RETURNING id", value, time.Now())
	if err := row.Err(); err != nil {
		s.log.Error("Error adding message", err)
		return uuid.Nil, err
	}
	s.countAll.Add(1)
	var id uuid.UUID
	if err := row.Scan(&id); err != nil {
		s.log.Error("Error scanning new message id", err)
		return uuid.Nil, err
	}
	return id, tx.Commit()
}

func (s *pgMessagesStore) GetOne(ctx context.Context, id uuid.UUID) (models.Message, error) {
	row := s.db.QueryRowContext(ctx, `SELECT id, value, processed, date_create, date_process
	FROM Messages WHERE id = $1`, id)
	if err := row.Err(); err != nil {
		s.log.Error("Error getting message", err)
		return models.Message{}, err
	}
	var msg models.Message
	if err := row.Scan(&msg.Id, &msg.Value, &msg.Processed, &msg.DateCreate, &msg.DateProcess); err != nil {
		s.log.Error("Error scanning message", err)
		return models.Message{}, err
	}
	return msg, nil
}

func (s *pgMessagesStore) MarkProcessed(ctx context.Context, id uuid.UUID) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		s.log.Error("Error marking message processed", err)
		return err
	}
	defer tx.Rollback()
	res, err := tx.ExecContext(ctx, `UPDATE Messages SET processed = true, date_process = $2 WHERE id = $1 AND processed = false`, id, time.Now())
	if err != nil {
		s.log.Error("Error marking message processed", err)
		return err
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		s.log.Error("can't update local processed messages count", err)
	} else if rowsAffected > 0 {
		s.countProcessed.Add(1)
	} else if rowsAffected == 0 {
		return sql.ErrNoRows
	}
	if err := tx.Commit(); err != nil {
		s.log.Error("error commiting processed message ", id)
		return err
	}
	return tx.Commit()
}

func (s *pgMessagesStore) Delete(ctx context.Context, id uuid.UUID) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		s.log.Error("error deleting message ", id.String(), " ", err.Error())
		return err
	}
	defer tx.Rollback()

	row := tx.QueryRowContext(ctx, "DELETE FROM Messages WHERE id = $1 RETURNING processed", id)
	var wasProcessed bool
	if err := row.Scan(&wasProcessed); err != nil {
		s.log.Error("error deleting message ", id.String(), " ", err.Error())
		return err
	}
	if err := tx.Commit(); err != nil {
		s.log.Error("error commiting deleted message ", id, " ", err.Error())
		return err
	}
	s.countAll.Add(-1)
	if wasProcessed {
		s.countProcessed.Add(-1)
	}
	return nil
}

func (s *pgMessagesStore) Stats(ctx context.Context) (models.Stats, error) {
	s.loadCountMu.Lock()
	if !s.countsLoaded {
		if err := s.loadCounts(ctx); err != nil {
			s.loadCountMu.Unlock()
			s.log.Error("error loading counts")
			return models.Stats{}, err
		}
	}
	s.loadCountMu.Unlock()

	return models.Stats{
		CountAll:       s.countAll.Load(),
		CountProcessed: s.countProcessed.Load(),
	}, nil
}

func (s *pgMessagesStore) loadCounts(ctx context.Context) error {
	countAllCh := make(chan int64)
	countProcessedCh := make(chan int64)
	group := errgroup.Group{}
	group.Go(func() error {
		row := s.db.QueryRowContext(ctx, `SELECT COUNT(id) FROM Messages`)
		var count int64
		if err := row.Scan(&count); err != nil {
			return err
		}
		go func() {
			countAllCh <- count
		}()
		return nil
	})
	group.Go(func() error {
		row := s.db.QueryRowContext(ctx, `SELECT COUNT(id) FROM Messages WHERE processed = true`)
		var count int64
		if err := row.Scan(&count); err != nil {
			return err
		}
		go func() {
			countProcessedCh <- count
		}()
		return nil
	})
	if err := group.Wait(); err != nil {
		return err
	}
	s.countAll.Store(<-countAllCh)
	s.countProcessed.Store(<-countProcessedCh)
	s.countsLoaded = true
	return nil
}
