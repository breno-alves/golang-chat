package services

import (
	"chatroom/internal/models"
	"chatroom/internal/repositories"
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type RoomService struct {
	RoomRepository *repositories.RoomRepository
}

func NewRoomService(db *gorm.DB, cache *redis.Client) *RoomService {
	return &RoomService{
		RoomRepository: repositories.NewRoomRepository(db, cache),
	}
}

func (rs *RoomService) CreateRoom(ctx context.Context) (*models.Room, error) {
	user := ctx.Value("user").(*models.User)
	roomName := fmt.Sprintf("%s's room.", user.Username)
	room, err := rs.RoomRepository.Create(roomName)
	if err != nil {
		return nil, err
	}
	return room, nil
}

func (rs *RoomService) FindAll(_ context.Context) (*[]models.Room, error) {
	rooms, err := rs.RoomRepository.FindAll()
	if err != nil {
		return nil, err
	}
	return rooms, nil
}

func (rs *RoomService) FindByID(_ context.Context, id uint) (*models.Room, error) {
	room, err := rs.RoomRepository.FindByID(id)
	if err != nil {
		return nil, err
	}
	return room, nil
}

func (rs *RoomService) UserJoinRoom(ctx context.Context) error {
	room := ctx.Value("room").(*models.Room)
	user := ctx.Value("user").(*models.User)
	token := ctx.Value("token").(string)
	_, err := rs.RoomRepository.SetRoomUserToken(room.Id, user.Id, token)
	if err != nil {
		return err
	}
	return nil
}

func (rs *RoomService) GetCurrentUserTokensInRoom(_ context.Context, roomId uint) ([]string, error) {
	tokens, err := rs.RoomRepository.GetUsersTokenInRoom(roomId)
	if err != nil {
		return nil, err
	}
	return tokens, nil
}

func (rs *RoomService) RemoveUserTokenInRoom(ctx context.Context) error {
	room, ok := ctx.Value("room").(*models.Room)
	if !ok {
		return fmt.Errorf("room not found in context")
	}

	user, ok := ctx.Value("user").(*models.User)
	if !ok {
		return fmt.Errorf("user not found in context")
	}

	err := rs.RoomRepository.RemoveUserTokenInRoom(room.Id, user.Id)
	if err != nil {
		return err
	}
	return nil
}
