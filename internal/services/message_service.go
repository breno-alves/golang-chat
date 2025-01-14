package services

import (
	"chatroom/internal/models"
	"chatroom/internal/pkg/broker"
	"chatroom/internal/repositories"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"log/slog"
	"regexp"
	"strings"
)

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
		slog.Error("could not find user by username")
		return nil, err
	}

	room, err := ms.roomRepository.FindByID(roomId)
	if err != nil {
		slog.Error("could not find room")
		return nil, err
	}

	message, err := ms.messageRepository.Create(user.Id, room.Id, content)
	if err != nil {
		return nil, err
	}

	// if message is a command will also create a request to bot API
	if strings.HasPrefix(message.Content, "/stock=") {
		stockCode := extractStockCode(message.Content)
		if stockCode != "" {
			if requestErr := ms.RequestStockPrice(message.RoomId, stockCode); requestErr != nil {
				slog.Error("could not request stock price")
			}
		}
	}
	return message, nil
}

func extractStockCode(input string) string {
	pattern := regexp.MustCompile(`/stock=\S+`)
	matching := pattern.FindAllString(input, -1)
	if len(matching) != 0 {
		firstCmd := matching[0]
		result := strings.Split(firstCmd, "=")[1]
		return result
	}
	return ""
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
		slog.Error("could not marshall")
		return err
	}
	channel.Ch <- bytes
	return nil
}

func (ms *MessageService) CreateMessageFromBot(response []byte) (*models.Message, error) {
	getStockPriceResponse := &GetStockPriceResponse{}
	if err := json.Unmarshal(response, getStockPriceResponse); err != nil {
		slog.Error("could not unmarshall")
		return nil, err
	}
	botUser, err := ms.userRepository.FindByUsername("bot")
	if err != nil {
		slog.Error("could not find bot")
		return nil, err
	}
	message, err := ms.messageRepository.Create(
		botUser.Id,
		getStockPriceResponse.RoomId,
		fmt.Sprintf("%s quote is $%v per share", strings.ToUpper(getStockPriceResponse.StockCode), getStockPriceResponse.Value))
	if err != nil {
		slog.Error("could not create message")
		return nil, err
	}
	return message, nil
}
