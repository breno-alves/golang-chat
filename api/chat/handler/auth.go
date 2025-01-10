package handler

import (
	"chatroom/internal/chat/service"
	"encoding/json"
	"gorm.io/gorm"
	"net/http"
)

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func Login(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	body := new(LoginRequest)
	err := json.NewDecoder(r.Body).Decode(body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	user, err := service.FindUserByUsername(db, body.Username)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if !service.CheckPasswordHash(user.Password, body.Password) {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	w.WriteHeader(http.StatusOK)
}
