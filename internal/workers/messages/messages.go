package messages

import (
	"context"
	"errors"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/YEgorLu/kafka-test/internal/logger"
	"github.com/YEgorLu/kafka-test/internal/mq/messages"
	serviceErrors "github.com/YEgorLu/kafka-test/internal/service/messages/errors"
	store "github.com/YEgorLu/kafka-test/internal/store/messages"
)

type MessagesWorker struct {
	readers          []*messages.MessagesConsumer
	s                store.MessagesStore
	log              logger.Logger
	numWorkers       int
	messages         chan messages.MessagesMessage
	closeCtx         context.Context
	cancelCtx        context.CancelFunc
	sysCloseNotifier sync.Once
}

func NewMessagesWorker(store store.MessagesStore, log logger.Logger) *MessagesWorker {
	numWorkers := 3
	readers := make([]*messages.MessagesConsumer, 0, numWorkers)
	for i := 0; i < numWorkers; i++ {
		readers = append(readers, messages.NewConsumer(log))
	}

	return &MessagesWorker{
		s:          store,
		log:        log,
		messages:   make(chan messages.MessagesMessage, numWorkers),
		numWorkers: numWorkers,
		readers:    readers,
	}
}

func (w *MessagesWorker) Start(ctx context.Context) {
	if w.cancelCtx != nil {
		w.cancelCtx()
	}
	w.closeCtx, w.cancelCtx = context.WithCancel(ctx)

	w.listenKafka(ctx)
	w.listenMessages(ctx)

	w.sysCloseNotifier.Do(func() {
		go func() {
			sigCh := make(chan os.Signal, 2)
			signal.Notify(sigCh, syscall.SIGTERM, syscall.SIGINT)
			for range sigCh {
				w.Stop()
			}
		}()
	})
}

func (w *MessagesWorker) Stop() {
	w.cancelCtx()
	w.closeCtx = nil
	w.cancelCtx = nil
	w.readers = []*messages.MessagesConsumer{}
}

func (w *MessagesWorker) listenKafka(ctx context.Context) {
	for _, reader := range w.readers {
		reader := reader
		go func() {
			defer reader.Close()
			for {
				msg, err := reader.ReadMessage(ctx)
				if err != nil {
					w.log.Error(err)
					if !errors.Is(err, context.Canceled) {
						continue
					}
				}
				w.log.Info("received message from kafka: ", msg.Value.Id.String())
				w.messages <- msg
			}
		}()
	}
}

func (w *MessagesWorker) listenMessages(ctx context.Context) {
	for range w.numWorkers {
		go func() {
			for {
				select {
				case <-ctx.Done():
					return
				case msg := <-w.messages:
					w.processMessage(ctx, msg)
				}
			}
		}()
	}
}

func (w *MessagesWorker) processMessage(ctx context.Context, msg messages.MessagesMessage) {
	if err := w.s.MarkProcessed(ctx, msg.Value.Id); err != nil && !errors.Is(err, serviceErrors.NonExistsError{Id: msg.Value.Id}) {
		w.log.Error("error processing msg ", msg.Value.Id)
	}
}
