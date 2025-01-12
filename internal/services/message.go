package services

import (
	"chatroom/internal/models"
	"fmt"
	"gorm.io/gorm"
	"log/slog"
)

const MaxMessagesToReturn = 50

func CreateMessage(db *gorm.DB, roomId uint, username, content string) (*models.Message, error) {
	user, err := FindUserByUsername(db, username)
	if err != nil {
		slog.Error("could not find user by username", username)
		return nil, err
	}
	room, err := FindRoomById(db, roomId)
	if err != nil {
		slog.Error("could not find room", roomId)
		return nil, err
	}
	message := models.NewMessage(user.Id, room.Id, content)
	err = db.Create(&message).Error
	if err != nil {
		return nil, err
	}
	return message, nil
}

func ListMessage(db *gorm.DB, roomId uint) (*[]models.Message, error) {
	messages := new([]models.Message)
	err := db.Order("created_at desc").Limit(MaxMessagesToReturn).Find(messages, "room_id = ?", roomId).Error
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return messages, nil
}
