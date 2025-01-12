package services

import (
	"chatroom/internal/models"
	"chatroom/internal/repositories"
	"context"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type UserService struct {
	UserRepository *repositories.UserRepository
}

func NewUserService(db *gorm.DB, cache *redis.Client) *UserService {
	return &UserService{
		UserRepository: repositories.NewUserRepository(db, cache),
	}
}

func (us *UserService) CreateUser(_ context.Context, username, password string) (*models.User, error) {
	user, err := us.UserRepository.Create(username, password)
	if err != nil {
		return nil, err
	}
	return user, err
}

func (us *UserService) FindUserByUsername(_ context.Context, username string) (*models.User, error) {
	user, err := us.UserRepository.FindByUsername(username)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (us *UserService) CheckPassword(ctx context.Context, username string, password string) (bool, error) {
	user, err := us.FindUserByUsername(ctx, username)
	if err != nil {
		return false, err
	}
	ok := user.Password == password
	return ok, nil
}

func (us *UserService) SetUserToken(_ context.Context, user *models.User) (string, error) {
	token, err := us.UserRepository.SetUserToken(user)
	if err != nil {
		return "", err
	}
	return token, nil
}

func (us *UserService) ValidateUserToken(_ context.Context, token string) (*models.User, error) {
	user, err := us.UserRepository.ValidateUserToken(token)
	if err != nil {
		return nil, err
	}
	return user, nil
}
