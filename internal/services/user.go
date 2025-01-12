package services

import (
	"chatroom/internal/models"
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"log/slog"
)

func CreateUser(db *gorm.DB, username, password string) (*models.User, error) {
	user := models.NewUser(username, password)
	err := db.Create(&user).Error
	if err != nil {
		return nil, err
	}
	return user, err
}

func FindUserByUsername(db *gorm.DB, username string) (*models.User, error) {
	user := &models.User{}
	err := db.First(&user, "username = ?", username).Error
	if err != nil {
		return nil, err
	}
	return user, err
}

func JoinRoom(ctx context.Context, _ *gorm.DB, cache *redis.Client) error {
	room := ctx.Value("room").(*models.Room)
	user := ctx.Value("user").(*models.User)
	token := ctx.Value("token").(string)
	err := cache.Set(ctx, fmt.Sprintf("room:%d:user:%d", room.Id, user.Id), token, 0).Err()
	if err != nil {
		slog.Error("failed to join room")
		return err
	}
	return nil
}
