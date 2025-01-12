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
	userRepository    *repositories.UserRepository
	roomRepository    *repositories.RoomRepository
	messageRepository *repositories.MessageRepository
}

func NewMessageService(db *gorm.DB, cache *redis.Client) *MessageService {
	return &MessageService{
		userRepository:    repositories.NewUserRepository(db, cache),
		roomRepository:    repositories.NewRoomRepository(db, cache),
		messageRepository: repositories.NewMessageRepository(db, cache),
	}
}

func (ms *MessageService) CreateMessage(_ context.Context, roomId uint, username, content string) (*models.Message, error) {
	user, err := ms.userRepository.FindByUsername(username)
	if err != nil {
		slog.Error("could not find user by username", username)
		return nil, err
	}

	room, err := ms.roomRepository.FindByID(roomId)
	if err != nil {
		slog.Error("could not find room", roomId)
		return nil, err
	}

	message, err := ms.messageRepository.Create(user.Id, room.Id, content)
	if err != nil {
		return nil, err
	}
	return message, nil
}

func (ms *MessageService) ListLastMessagesFromRoom(_ context.Context, roomId uint) (*[]models.Message, error) {
	messages, err := ms.messageRepository.FindLastMessagesByRoomId(roomId)
	if err != nil {
		return nil, err
	}
	return messages, nil
}
