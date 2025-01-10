package handler

import (
	"chatroom/internal/chat/service"
	"encoding/json"
	"gorm.io/gorm"
	"net/http"
)

type SignUpRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func SignUp(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
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
	if err := json.NewEncoder(w).Encode(user); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
}
