package models

import (
	"golang.org/x/crypto/bcrypt"
	"log/slog"
	"time"
)

type User struct {
	Id        uint      `gorm:"primaryKey" json:"id,omitempty"`
	Username  string    `gorm:"unique" json:"username,omitempty"`
	Password  string    `gorm:"not null" json:"-"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
}

func NewUser(username, password string) (*User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	if err != nil {
		slog.Error("error hashing password", err.Error())
		return nil, err
	}
	return &User{
		Username:  username,
		Password:  string(hashedPassword),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}, nil
}
