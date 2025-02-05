package chat

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
)

func (h *Handler) CreateRoom(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	token := r.URL.Query().Get("token")
	ctx = context.WithValue(ctx, "token", token)

	user, err := h.UserService.FindUserByToken(ctx, token)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	ctx = context.WithValue(ctx, "user", user)

	room, err := h.RoomService.CreateRoom(ctx)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	err = json.NewEncoder(w).Encode(room)
	if err != nil {
		slog.Error(fmt.Sprintf("could not encode rooms %s", err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (h *Handler) ListRooms(w http.ResponseWriter, r *http.Request) {
	slog.Debug("attempting to list rooms")
	ctx := r.Context()

	roomsList, err := h.RoomService.FindAll(ctx)
	if err != nil {
		slog.Error(fmt.Sprintf("could not list rooms: %s", err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(roomsList)
	if err != nil {
		slog.Error(fmt.Sprintf("could not encode rooms: %s", err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (h *Handler) LeaveRoom(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	token := r.URL.Query().Get("token")
	user, err := h.UserService.FindUserByToken(ctx, token)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	ctx = context.WithValue(ctx, "user", user)

	roomId, err := strconv.Atoi(r.URL.Query().Get("room_id"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	room, err := h.RoomService.FindByID(ctx, uint(roomId))
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	ctx = context.WithValue(ctx, "room", room)
	err = h.RoomService.RemoveUserTokenInRoom(ctx)
	if err != nil {
		slog.Error(fmt.Sprintf("could not remove user token in room: %s", err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
