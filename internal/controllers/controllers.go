package controllers

import (
	"context"
	"net/http"
	"sync"

	"github.com/YEgorLu/kafka-test/internal/config"
	"github.com/YEgorLu/kafka-test/internal/controllers/messages"
	"github.com/YEgorLu/kafka-test/internal/db"
	"github.com/YEgorLu/kafka-test/internal/logger"
	messagesMq "github.com/YEgorLu/kafka-test/internal/mq/messages"
	messagesService "github.com/YEgorLu/kafka-test/internal/service/messages"
	messagesStore "github.com/YEgorLu/kafka-test/internal/store/messages"
	messagesWorker "github.com/YEgorLu/kafka-test/internal/workers/messages"
)

type Controller interface {
	RegisterRoute(*http.ServeMux)
}

var onceInit sync.Once
var controllers []Controller

func GetRoutes() (*http.ServeMux, error) {
	onceInit.Do(func() {
		initControllers()
	})
	router := http.NewServeMux()
	for _, controller := range controllers {
		controller.RegisterRoute(router)
	}
	return router, nil
}

func initControllers() {
	log := logger.Get()
	log.Info("config ", config.App, config.DB, config.Mq)
	db := db.GetConnection(log)

	messagesStore := messagesStore.NewStore(db, log)

	messagesProducer := messagesMq.NewMessagesProducer(log)
	errMessagesProducer := messagesMq.NewErrorMessagesProducer(log)
	messagesService := messagesService.NewMessagesService(messagesStore, messagesProducer, errMessagesProducer, log)

	controllers = []Controller{
		messages.NewController(messagesService, log),
	}
	messageWorker := messagesWorker.NewMessagesWorker(messagesStore, log)
	messageWorker.Start(context.Background())
}
