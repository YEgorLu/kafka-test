package converters

import (
	"github.com/YEgorLu/kafka-test/internal/service/messages/models"
	storeModels "github.com/YEgorLu/kafka-test/internal/store/messages/models"
)

func StatsToSv(storeModel storeModels.Stats) models.Stats {
	return models.Stats{
		CountAll:       storeModel.CountAll,
		CountProcessed: storeModel.CountProcessed,
	}
}
