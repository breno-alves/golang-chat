package service

import (
	"chatroom/api/chat/model"
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"log/slog"
)

func CreateUser(db *gorm.DB, username, password string) (*model.User, error) {
	user := model.NewUser(username, password)
	err := db.Create(&user).Error
	if err != nil {
		return nil, err
	}
	return user, err
}

func FindUserByUsername(db *gorm.DB, username string) (*model.User, error) {
	user := &model.User{}
	err := db.First(&user, "username = ?", username).Error
	if err != nil {
		return nil, err
	}
	return user, err
}

func JoinRoom(ctx context.Context, _ *gorm.DB, cache *redis.Client) error {
	room := ctx.Value("room").(*model.Room)
	user := ctx.Value("user").(*model.User)
	token := ctx.Value("token").(string)
	err := cache.Set(ctx, fmt.Sprintf("room:%d:user:%d", room.Id, user.Id), token, 0).Err()
	if err != nil {
		slog.Error("failed to join room")
		return err
	}
	return nil
}
