package service

import (
	"chatroom/api/chat/model"
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"time"
)

func CheckPasswordHash(hash, password string) bool {
	return hash == password
}

func ValidateUserToken(ctx context.Context, cache *redis.Client, token string) (*model.User, error) {
	key := fmt.Sprintf("user:%s", token)
	data, err := cache.Get(ctx, key).Result()
	if err != nil {
		return nil, err
	}
	user := &model.User{}
	err = json.Unmarshal([]byte(data), user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func SetUserToken(ctx context.Context, cache *redis.Client, user *model.User) (string, error) {
	payload, err := json.Marshal(user)
	if err != nil {
		return "", err
	}
	hash := uuid.New().String()
	key := fmt.Sprintf("user:%s", hash)
	err = cache.Set(ctx, key, payload, time.Hour*24).Err()
	if err != nil {
		return "", err
	}
	return hash, nil
}
