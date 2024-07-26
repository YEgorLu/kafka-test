package converters

import (
	"github.com/YEgorLu/kafka-test/internal/controllers/messages/models"
	serviceModels "github.com/YEgorLu/kafka-test/internal/service/messages/models"
)

func StatsFromSv(serviceModel serviceModels.Stats) models.StatsResponse {
	return models.StatsResponse{
		CountAll:       serviceModel.CountAll,
		CountProcessed: serviceModel.CountProcessed,
	}
}
