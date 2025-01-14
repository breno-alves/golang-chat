package chat

import (
	"encoding/json"
	"errors"
	"fmt"
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

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	body := new(LoginRequest)
	err := json.NewDecoder(r.Body).Decode(body)
	if err != nil {
		slog.Error(fmt.Sprintf("failed to decode body %s", err.Error()))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	user, err := h.UserService.FindUserByUsername(ctx, body.Username)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			slog.Error("user not found")
			w.WriteHeader(http.StatusNotFound)
			return
		}
		slog.Error(fmt.Sprintf("failed to find user %s", err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	valid, err := h.UserService.CheckPassword(ctx, user.Username, body.Password)
	if err != nil {
		slog.Error(fmt.Sprintf("failed to check password %s", err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
	}

	if !valid {
		slog.Warn("password does not match")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	hash, err := h.UserService.SetUserToken(ctx, user)
	if err != nil {
		slog.Error(fmt.Sprintf("failed to set token %s", err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	err = json.NewEncoder(w).Encode(&LoginResponse{
		Token: hash,
	})
	if err != nil {
		slog.Error(fmt.Sprintf("failed to encode response %s", err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	slog.Debug("successfully logged in")
}
