package services

import (
	"chatroom/internal/models"
	"chatroom/internal/repositories"
	"context"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"log/slog"
)

const MaxMessagesToReturn = 50

type MessageService struct {
	UserRepository    *repositories.UserRepository
	RoomRepository    *repositories.RoomRepository
	MessageRepository *repositories.MessageRepository
}

func NewMessageService(db *gorm.DB, cache *redis.Client) *MessageService {
	return &MessageService{
		UserRepository:    repositories.NewUserRepository(db, cache),
		RoomRepository:    repositories.NewRoomRepository(db, cache),
		MessageRepository: repositories.NewMessageRepository(db, cache),
	}
}

func (ms *MessageService) CreateMessage(_ context.Context, roomId uint, username, content string) (*models.Message, error) {
	user, err := ms.UserRepository.FindByUsername(username)
	if err != nil {
		slog.Error("could not find user by username", username)
		return nil, err
	}

	room, err := ms.RoomRepository.FindByID(roomId)
	if err != nil {
		slog.Error("could not find room", roomId)
		return nil, err
	}

	message, err := ms.MessageRepository.Create(user.Id, room.Id, content)
	if err != nil {
		return nil, err
	}
	return message, nil
}

//func ListMessage(db *gorm.DB, roomId uint) (*[]models.Message, error) {
//	messages := new([]models.Message)
//	err := db.Order("created_at desc").Limit(MaxMessagesToReturn).Find(messages, "room_id = ?", roomId).Error
//	if err != nil {
//		fmt.Println(err)
//		return nil, err
//	}
//	return messages, nil
//}
