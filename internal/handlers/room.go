package handlers

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

func (h *Handler) ListRooms(w http.ResponseWriter, r *http.Request) {
	slog.Debug("attempting to list rooms")
	ctx := r.Context()

	roomsList, err := h.roomService.FindAll(ctx)
	if err != nil {
		slog.Error("could not list rooms: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(roomsList)
	if err != nil {
		slog.Error("could not encode rooms: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
