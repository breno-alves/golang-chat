package chat

import (
	"encoding/json"
	"errors"
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

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	body := new(LoginRequest)
	err := json.NewDecoder(r.Body).Decode(body)
	if err != nil {
		slog.Error("failed to decode body", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	user, err := h.UserService.FindUserByUsername(ctx, body.Username)
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

	valid, err := h.UserService.CheckPassword(ctx, user.Username, body.Password)
	if err != nil {
		slog.Error("failed to check password", err)
		w.WriteHeader(http.StatusInternalServerError)
	}

	if !valid {
		slog.Warn("password does not match", err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	hash, err := h.UserService.SetUserToken(ctx, user)
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
