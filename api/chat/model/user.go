package model

import (
	"database/sql"
	"time"
)

type User struct {
	Id         uint          `gorm:"primaryKey" json:"id,omitempty"`
	Username   string        `gorm:"unique" json:"username,omitempty"`
	Password   string        `gorm:"not null" json:"password,omitempty"`
	RoomRefer  sql.NullInt32 `gorm:"null" json:"active_room,omitempty"`
	ActiveRoom Room          `json:"active_room,omitempty" gorm:"foreignKey:RoomRefer"`
	CreateAt   time.Time     `json:"create_at,omitempty"`
	UpdateAt   time.Time     `json:"update_at,omitempty"`
}

func NewUser(username, password string) *User {
	return &User{
		Username: username,
		Password: password,
		CreateAt: time.Now(),
		UpdateAt: time.Now(),
	}
}
