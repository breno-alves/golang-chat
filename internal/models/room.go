package models

import (
	"time"
)

type Room struct {
	Id        uint      `gorm:"primaryKey" json:"id"`
	Title     string    `json:"title"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
}

func NewRoom(title string) *Room {
	return &Room{
		Title:     title,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}
