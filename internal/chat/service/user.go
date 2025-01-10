package service

import (
	"chatroom/api/chat/model"
	"gorm.io/gorm"
)

func CreateUser(db *gorm.DB, username, password string) (*model.User, error) {
	user := model.NewUser(username, password)
	err := db.Create(&user).Error
	if err != nil {
		return nil, err
	}
	return user, err
}

func FindUserByUsername(db *gorm.DB, username string) (*model.User, error) {
	user := &model.User{}
	err := db.First(&user, "username = ?", username).Error
	if err != nil {
		return nil, err
	}
	return user, err
}
