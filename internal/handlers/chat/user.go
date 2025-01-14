package chat

import (
	"encoding/json"
	"fmt"
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
		slog.Error(fmt.Sprintf("failed to decode body %s", err.Error()))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	user, err := h.UserService.CreateUser(ctx, body.Username, body.Password)
	if err != nil {
		slog.Error(fmt.Sprintf("failed to create user %s", err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(user)
	if err != nil {
		slog.Error(fmt.Sprintf("failed to encode response %s", err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	slog.Debug("successfully created user")
}
