package repositories

import (
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type AuthRepository struct {
	db    *gorm.DB
	cache *redis.Client
}

func NewAuthRepository(db *gorm.DB, cache *redis.Client) *AuthRepository {
	return &AuthRepository{
		db:    db,
		cache: cache,
	}
}
