package handler

import (
	"chatroom/internal/chat/service"
	"encoding/json"
	"fmt"
	"gorm.io/gorm"
	"log/slog"
	"net/http"
)

type JoinRoomRequest struct {
	Username string `json:"username"`
	RoomID   uint   `json:"room_id"`
}

func CreateRoom(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	room, err := service.CreateRoom(db)
	if err != nil {
		slog.Error(err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if err := json.NewEncoder(w).Encode(room); err != nil {
		slog.Error(err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func ListRooms(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	rooms, err := service.ListRooms(db)
	if err != nil {
		slog.Error(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(rooms); err != nil {
		slog.Error(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func JoinRoom(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	body := new(JoinRoomRequest)
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		slog.Error(err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	user, err := service.FindUserByUsername(db, body.Username)
	if err != nil {
		slog.Error(err.Error())
		w.WriteHeader(http.StatusNotFound)
		return
	}

	fmt.Println(user)

	room, err := service.FindRoomById(db, body.RoomID)
	if err != nil {
		slog.Error(err.Error())
		w.WriteHeader(http.StatusNotFound)
		return
	}
	fmt.Println(room)
}
