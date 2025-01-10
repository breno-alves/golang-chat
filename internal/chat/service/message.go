package service

import (
	"chatroom/api/chat/model"
	"fmt"
	"gorm.io/gorm"
	"log/slog"
)

func CreateMessage(db *gorm.DB, roomId uint, username, content string) (*model.Message, error) {
	// find owner
	user, err := FindUserByUsername(db, username)
	if err != nil {
		slog.Debug("could not find user by username", username)
		return nil, err
	}
	fmt.Println(user)
	// find room
	room, err := FindRoomById(db, roomId)
	if err != nil {
		slog.Debug("could not find room", roomId)
		return nil, err
	}
	fmt.Println(room)
	// create message
	return nil, nil
}
