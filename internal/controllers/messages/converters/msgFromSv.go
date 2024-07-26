package converters

import (
	"github.com/YEgorLu/kafka-test/internal/controllers/messages/models"
	serviceModels "github.com/YEgorLu/kafka-test/internal/service/messages/models"
)

func MsgFromSv(serviceModel serviceModels.Message) models.GetOneResponse {
	return models.GetOneResponse{
		Id:          serviceModel.Id,
		Value:       serviceModel.Value,
		Processed:   serviceModel.Processed,
		DateCreate:  serviceModel.DateCreate,
		DateProcess: serviceModel.DateProcess,
	}
}
