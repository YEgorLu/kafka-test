package messages

import (
	"net/http"

	"github.com/YEgorLu/kafka-test/internal/logger"
	"github.com/YEgorLu/kafka-test/internal/service/messages"
	"github.com/YEgorLu/kafka-test/internal/util"
)

type MessagesController struct {
	log logger.Logger
	s   messages.MessageService
}

func NewController(s messages.MessageService, log logger.Logger) *MessagesController {
	return &MessagesController{
		log: log,
		s:   s,
	}
}

func (c *MessagesController) RegisterRoute(router *http.ServeMux) {
	p := util.Rpm("messages")
	router.HandleFunc(p(http.MethodPost), c.Add)
	router.HandleFunc(p(http.MethodGet, "stats"), c.Stats)
	router.HandleFunc(p(http.MethodGet, "{id}"), c.GetOne)
}
