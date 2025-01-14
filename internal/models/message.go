package models

import "time"

type Message struct {
	Id        uint      `gorm:"primaryKey" json:"id"`
	Content   string    `json:"content"`
	OwnerId   uint      `gorm:"null" json:"owner_id"`
	Owner     User      `gorm:"foreignKey:OwnerId" json:"owner"`
	RoomId    uint      `json:"room_id"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
}

func NewMessage(ownerId uint, roomId uint, content string) *Message {
	return &Message{
		OwnerId:   ownerId,
		Content:   content,
		RoomId:    roomId,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}
