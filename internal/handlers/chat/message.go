package chat

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"
)

func (h *Handler) ListMessages(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	roomId, err := strconv.Atoi(r.URL.Query().Get("room_id"))
	if err != nil {
		slog.Error("room_id is required")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	messages, err := h.MessageService.ListLastMessagesFromRoom(ctx, uint(roomId))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	err = json.NewEncoder(w).Encode(messages)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
