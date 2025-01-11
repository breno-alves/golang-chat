package handler

import (
	"chatroom/internal/chat/service"
	"context"
	"encoding/json"
	"errors"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"log/slog"
	"net/http"
)

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

// Login steps
// - Check valid login
// - Create new token (uuid)
// - Add token to redis with user data
// - Return token or 403

func Login(ctx context.Context, db *gorm.DB, cache *redis.Client, w http.ResponseWriter, r *http.Request) {
	slog.Debug("attempting to login user")

	body := new(LoginRequest)
	err := json.NewDecoder(r.Body).Decode(body)
	if err != nil {
		slog.Error("failed to decode body", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	user, err := service.FindUserByUsername(db, body.Username)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			slog.Error("user not found", err)
			w.WriteHeader(http.StatusNotFound)
			return
		}
		slog.Error("failed to find user", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if !service.CheckPasswordHash(user.Password, body.Password) {
		slog.Warn("password does not match", err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	hash, err := service.SetUserToken(ctx, cache, user)
	if err != nil {
		slog.Error("failed to set token", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	err = json.NewEncoder(w).Encode(&LoginResponse{
		Token: hash,
	})
	if err != nil {
		slog.Error("failed to encode response", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	slog.Debug("successfully logged in")
}
