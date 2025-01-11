package handler

import (
	"chatroom/internal/chat/service"
	"context"
	"encoding/json"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"log/slog"
	"net/http"
)

func ListRooms(_ context.Context, db *gorm.DB, _ *redis.Client, w http.ResponseWriter, _ *http.Request) {
	slog.Debug("attempting to list rooms")
	roomsList, err := service.ListRooms(db)
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
