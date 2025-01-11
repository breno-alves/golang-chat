package model

import (
	"time"
)

type Room struct {
	Id        uint      `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
}

func NewRoom() *Room {
	return &Room{
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}
