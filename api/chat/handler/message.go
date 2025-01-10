package handler

import (
	"chatroom/internal/chat/service"
	"encoding/json"
	"github.com/gorilla/mux"
	"gorm.io/gorm"
	"net/http"
	"strconv"
)

type CreateMessageRequest struct {
	RoomId   uint   `json:"room_id"`
	Username string `json:"username"`
	Content  string `json:"content"`
}

func ListMessage(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	roomId, err := strconv.Atoi(mux.Vars(r)["room_id"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	messages, err := service.ListMessage(db, uint(roomId))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(messages)
}

func CreateMessage(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	body := new(CreateMessageRequest)
	err := json.NewDecoder(r.Body).Decode(body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	message, err := service.CreateMessage(db, body.RoomId, body.Username, body.Content)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	json.NewEncoder(w).Encode(message)
}
