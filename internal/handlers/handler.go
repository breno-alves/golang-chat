package handlers

import (
	"chatroom/internal/services"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type Handler struct {
	userService    *services.UserService
	roomService    *services.RoomService
	messageService *services.MessageService
	authService    *services.AuthService
}

func NewHandler(db *gorm.DB, cache *redis.Client) *Handler {
	return &Handler{
		userService:    services.NewUserService(db, cache),
		roomService:    services.NewRoomService(db, cache),
		messageService: services.NewMessageService(db, cache),
		//authService:    services.NewAuthService(db, cache),
	}
}
