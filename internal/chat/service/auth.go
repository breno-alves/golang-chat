package service

import (
	"context"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"time"
)

func CheckPasswordHash(hash, password string) bool {
	return hash == password
}

func SetToken(ctx context.Context, cache *redis.Client, username string) (string, error) {
	hash := uuid.New().String()
	err := cache.Set(ctx, hash, username, time.Hour*24).Err()
	if err != nil {
		return "", err
	}
	return hash, nil
}
