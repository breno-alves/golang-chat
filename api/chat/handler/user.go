package handler

import (
	"chatroom/internal/chat/service"
	"context"
	"encoding/json"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"net/http"
)

type SignUpRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type SignUpResponse struct {
	Token string `json:"token"`
}

func SignUp(ctx context.Context, db *gorm.DB, cache *redis.Client, w http.ResponseWriter, r *http.Request) {
	body := new(SignUpRequest)
	err := json.NewDecoder(r.Body).Decode(body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	user, err := service.CreateUser(db, body.Username, body.Password)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	token, err := service.SetToken(ctx, cache, user.Username)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(&SignUpResponse{Token: token}); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
}
