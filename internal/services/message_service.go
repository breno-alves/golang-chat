package services

import (
	"chatroom/internal/models"
	"chatroom/internal/pkg/broker"
	"chatroom/internal/repositories"
	"context"
	"encoding/json"
	"errors"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"log/slog"
)

const MaxMessagesToReturn = 50

type MessageService struct {
	userRepository    *repositories.UserRepository
	roomRepository    *repositories.RoomRepository
	messageRepository *repositories.MessageRepository
	broker            *broker.Broker
}

func NewMessageService(db *gorm.DB, cache *redis.Client, broker *broker.Broker) *MessageService {
	return &MessageService{
		userRepository:    repositories.NewUserRepository(db, cache),
		roomRepository:    repositories.NewRoomRepository(db, cache),
		messageRepository: repositories.NewMessageRepository(db, cache),
		broker:            broker,
	}
}

func (ms *MessageService) CreateMessage(_ context.Context, roomId uint, username, content string) (*models.Message, error) {
	user, err := ms.userRepository.FindByUsername(username)
	if err != nil {
		slog.Error("could not find user by username", username)
		return nil, err
	}

	room, err := ms.roomRepository.FindByID(roomId)
	if err != nil {
		slog.Error("could not find room", roomId)
		return nil, err
	}

	message, err := ms.messageRepository.Create(user.Id, room.Id, content)
	if err != nil {
		return nil, err
	}
	go func() {
		err := ms.RequestStockPrice(message.RoomId, "aapl.us")
		if err != nil {
			slog.Error("could not request stock price", message.RoomId)
		}
	}()
	return message, nil
}

func (ms *MessageService) ListLastMessagesFromRoom(_ context.Context, roomId uint) (*[]models.Message, error) {
	messages, err := ms.messageRepository.FindLastMessagesByRoomId(roomId)
	if err != nil {
		return nil, err
	}
	return messages, nil
}

type RequestStockPricePayload struct {
	RoomId    uint   `json:"room_id"`
	StockCode string `json:"stock_code"`
}

func (ms *MessageService) RequestStockPrice(roomId uint, stockCode string) error {
	slog.Debug("sending asdasuidasduasd")
	channel, ok := ms.broker.Channels["BOT_STOCKS"]
	if !ok {
		slog.Error("could not find broker channel")
		return errors.New("could not find broker channel")
	}
	bytes, err := json.Marshal(&RequestStockPricePayload{
		RoomId:    roomId,
		StockCode: stockCode,
	})
	if err != nil {
		slog.Error("could not marshal seila")
		return err
	}
	channel.Ch <- bytes
	return nil
}
