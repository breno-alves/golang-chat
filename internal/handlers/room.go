package handlers

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
)

func (h *Handler) CreateRoom(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	token := r.URL.Query().Get("token")
	ctx = context.WithValue(ctx, "token", token)

	user, err := h.userService.FindUserByToken(ctx, token)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	ctx = context.WithValue(ctx, "user", user)

	room, err := h.roomService.CreateRoom(ctx)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	err = json.NewEncoder(w).Encode(room)
	if err != nil {
		slog.Error("could not encode rooms", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

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
