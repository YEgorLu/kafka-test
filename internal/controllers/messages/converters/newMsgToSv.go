package converters

import (
	"github.com/YEgorLu/kafka-test/internal/controllers/messages/models"
	serviceModels "github.com/YEgorLu/kafka-test/internal/service/messages/models"
)

func NewMsgToSv(reqMsg models.NewMessage) serviceModels.NewMessage {
	return serviceModels.NewMessage{
		Value: reqMsg.Value,
	}
}
