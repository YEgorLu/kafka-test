package converters

import (
	"github.com/YEgorLu/kafka-test/internal/service/messages/models"
	storeModels "github.com/YEgorLu/kafka-test/internal/store/messages/models"
)

func MsgToSv(storeMsg storeModels.Message) models.Message {
	return models.Message{
		Id:          storeMsg.Id,
		Value:       storeMsg.Value,
		Processed:   storeMsg.Processed,
		DateCreate:  storeMsg.DateCreate,
		DateProcess: storeMsg.DateProcess,
	}
}
