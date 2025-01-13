package chat

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

type SignUpRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (h *Handler) SignUp(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	body := &SignUpRequest{}
	err := json.NewDecoder(r.Body).Decode(body)
	if err != nil {
		slog.Error("failed to decode body", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	user, err := h.userService.CreateUser(ctx, body.Username, body.Password)
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
