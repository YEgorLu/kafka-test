package messages

import (
	"encoding/json"
	"net/http"

	"github.com/YEgorLu/kafka-test/internal/controllers/messages/converters"
	"github.com/YEgorLu/kafka-test/internal/util"
)

// Get One godoc
// @Tags messages
// @Summary Returns message
// @Produce json
// @Param id path string true "Message ID"
// @Success 200 {object} models.GetOneResponse
// @Router /messages/{id} [get]
func (c *MessagesController) GetOne(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	idPath := r.PathValue("id")
	id, err := util.ParseUUID(idPath)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	serviceMsg, err := c.s.GetOne(ctx, id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if err := json.NewEncoder(w).Encode(converters.MsgFromSv(serviceMsg)); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
