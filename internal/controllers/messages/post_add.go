package messages

import (
	"encoding/json"
	"net/http"

	"github.com/YEgorLu/kafka-test/internal/controllers/messages/converters"
	"github.com/YEgorLu/kafka-test/internal/controllers/messages/models"
)

// Add godoc
// @Tags messages
// @Summary Creates message and puts it in queue for later processing
// @Produce json
// @Param value body string true "Your message" SchemaExample({"value": "your message"})
// @Success 200 {object} models.AddResponse
// @Router /messages [post]
func (c *MessagesController) Add(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var body models.NewMessage
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		c.log.Error(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	id, err := c.s.Add(ctx, converters.NewMsgToSv(body))
	if err != nil {
		c.log.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if err := json.NewEncoder(w).Encode(models.AddResponse{Id: id}); err != nil {
		c.log.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}
