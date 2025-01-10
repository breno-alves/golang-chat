package model

import (
	"time"
)

type Room struct {
	Id       uint      `json:"id"`
	CreateAt time.Time `json:"create_at,omitempty"`
	UpdateAt time.Time `json:"update_at,omitempty"`
}

func NewRoom() *Room {
	return &Room{}
}
