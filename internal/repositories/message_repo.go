package repositories

import (
	"chatroom/internal/models"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type MessageRepository struct {
	db    *gorm.DB
	cache *redis.Client
}

func NewMessageRepository(db *gorm.DB, cache *redis.Client) *MessageRepository {
	return &MessageRepository{
		db:    db,
		cache: cache,
	}
}

func (mr *MessageRepository) Create(userId, roomId uint, content string) (*models.Message, error) {
	message := models.NewMessage(userId, roomId, content)
	err := mr.db.Create(message).Error
	if err != nil {
		return nil, err
	}
	return message, nil
}

func (mr *MessageRepository) FindLastMessagesByRoomId(roomId uint) (*[]models.Message, error) {
	messages := &[]models.Message{}
	err := mr.db.Limit(50).Find(&messages, "room_id = ?", roomId).Error
	if err != nil {
		return nil, err
	}
	return messages, nil
}
