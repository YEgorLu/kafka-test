package messages

import (
	"encoding/json"
	"net/http"

	"github.com/YEgorLu/kafka-test/internal/controllers/messages/converters"
)

// Stats godoc
// @Tags messages
// @Summary Returns overall and processed messages count
// @Produce json
// @Success 200 {object} models.StatsResponse
// @Router /messages/stats [get]
func (c *MessagesController) Stats(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	stats, err := c.s.Stats(ctx)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if err := json.NewEncoder(w).Encode(converters.StatsFromSv(stats)); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
