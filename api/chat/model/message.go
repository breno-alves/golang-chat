package model

import "time"

type Message struct {
	Id       uint      `json:"id"`
	Content  string    `json:"content"`
	Owner    uint      `json:"owner"`
	CreateAt time.Time `json:"create_at,omitempty"`
	UpdateAt time.Time `json:"update_at,omitempty"`
}

func NewMessage() *Message {
	return &Message{}
}
