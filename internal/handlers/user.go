package handlers

import (
	"chatroom/internal/services"
	"context"
	"encoding/json"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"log/slog"
	"net/http"
)

type SignUpRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func SignUp(_ context.Context, db *gorm.DB, _ *redis.Client, w http.ResponseWriter, r *http.Request) {
	slog.Debug("attempting to sign up user")

	body := new(SignUpRequest)
	err := json.NewDecoder(r.Body).Decode(body)
	if err != nil {
		slog.Error("failed to decode body", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	user, err := services.CreateUser(db, body.Username, body.Password)
	if err != nil {
		slog.Error("failed to create user", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(user)
	if err != nil {
		slog.Error("failed to encode response", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	slog.Debug("successfully created user")
}
