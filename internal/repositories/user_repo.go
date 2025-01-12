package repositories

import (
	"chatroom/internal/models"
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"time"
)

type UserRepository struct {
	db    *gorm.DB
	cache *redis.Client
}

func NewUserRepository(db *gorm.DB, cache *redis.Client) *UserRepository {
	return &UserRepository{
		db:    db,
		cache: cache,
	}
}

func (u *UserRepository) Create(username, password string) (*models.User, error) {
	user, err := models.NewUser(username, password)
	if err != nil {
		return nil, err
	}
	err = u.db.Create(user).Error
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (u *UserRepository) FindByUsername(username string) (*models.User, error) {
	user := &models.User{}
	err := u.db.Where("username = ?", username).First(user).Error
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (u *UserRepository) SetUserToken(user *models.User) (string, error) {
	ctx := context.Background()
	payload, err := json.Marshal(user)
	if err != nil {
		return "", err
	}
	hash := uuid.New().String()
	key := fmt.Sprintf("user:%s", hash)
	err = u.cache.Set(ctx, key, payload, time.Hour*24).Err()
	if err != nil {
		return "", err
	}
	return hash, nil
}

func (u *UserRepository) ValidateUserToken(token string) (*models.User, error) {
	ctx := context.Background()
	key := fmt.Sprintf("user:%s", token)
	data, err := u.cache.Get(ctx, key).Result()
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

func (u *UserRepository) FindUserByToken(token string) (*models.User, error) {
	ctx := context.Background()
	key := fmt.Sprintf("user:%s", token)
	data, err := u.cache.Get(ctx, key).Result()
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
