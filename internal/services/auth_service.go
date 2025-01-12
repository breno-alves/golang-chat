package services

import (
	"chatroom/internal/models"
	"chatroom/internal/repositories"
	"context"
	"encoding/json"
	"fmt"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type AuthService struct {
	AuthRepository *repositories.AuthRepository
}

func NewAuthService(db *gorm.DB, cache *redis.Client) *AuthService {
	return &AuthService{
		AuthRepository: repositories.NewAuthRepository(db, cache),
	}
}

func ValidateUserToken(ctx context.Context, cache *redis.Client, token string) (*models.User, error) {
	key := fmt.Sprintf("user:%s", token)
	data, err := cache.Get(ctx, key).Result()
	if err != nil {
		return nil, err
	}
	user := &models.User{}
	err = json.Unmarshal([]byte(data), user)
	if err != nil {
		return nil, err
	}
	return user, nil
}
