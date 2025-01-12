package repositories

import (
	"chatroom/internal/models"
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"log/slog"
)

type RoomRepository struct {
	db    *gorm.DB
	cache *redis.Client
}

func NewRoomRepository(db *gorm.DB, cache *redis.Client) *RoomRepository {
	return &RoomRepository{
		db:    db,
		cache: cache,
	}
}

func (rr *RoomRepository) Create(title string) (*models.Room, error) {
	room := models.NewRoom(title)
	err := rr.db.Create(&room).Error
	if err != nil {
		return nil, err
	}
	return room, nil
}

func (rr *RoomRepository) FindAll() (*[]models.Room, error) {
	rooms := &[]models.Room{}
	err := rr.db.Find(rooms).Error
	if err != nil {
		return nil, err
	}
	return rooms, nil
}

func (rr *RoomRepository) FindByID(id uint) (*models.Room, error) {
	room := &models.Room{}
	err := rr.db.Where("id = ?", id).Find(room).Error
	if err != nil {
		return nil, err
	}
	return room, nil
}

func (rr *RoomRepository) SetRoomUserToken(roomId, userId uint, token string) (string, error) {
	ctx := context.Background()
	key := fmt.Sprintf("room:%d:user:%d", roomId, userId)
	err := rr.cache.Set(ctx, key, token, 0).Err()
	if err != nil {
		slog.Error("failed to join room")
		return "", err
	}
	return key, nil
}

func (rr *RoomRepository) GetUsersTokenInRoom(roomId uint) ([]string, error) {
	ctx := context.Background()

	var result []string

	iter := rr.cache.Scan(ctx, 0, fmt.Sprintf("room:%d:user:*", roomId), 0).Iterator()
	for iter.Next(ctx) {
		k := iter.Val()
		token, _ := rr.cache.Get(ctx, k).Result()
		result = append(result, token)
	}
	if err := iter.Err(); err != nil {
		return nil, err
	}
	return result, nil
}
