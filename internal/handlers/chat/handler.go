package chat

import (
	"chatroom/internal/pkg/broker"
	"chatroom/internal/services"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type Handler struct {
	UserService    *services.UserService
	RoomService    *services.RoomService
	MessageService *services.MessageService
}

func NewHandler(db *gorm.DB, cache *redis.Client, broker *broker.Broker) *Handler {
	return &Handler{
		UserService:    services.NewUserService(db, cache),
		RoomService:    services.NewRoomService(db, cache),
		MessageService: services.NewMessageService(db, cache, broker),
	}
}
