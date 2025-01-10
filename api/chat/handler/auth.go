package handler

import (
	"chatroom/internal/chat/service"
	"context"
	"encoding/json"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"net/http"
)

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func Login(_ context.Context, db *gorm.DB, _ *redis.Client, w http.ResponseWriter, r *http.Request) {
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

	json.NewEncoder(w).Encode(user)
}
