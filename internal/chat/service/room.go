package service

import (
	"chatroom/api/chat/model"
	"gorm.io/gorm"
)

func CreateRoom(db *gorm.DB) (*model.Room, error) {
	room := model.NewRoom()
	err := db.Create(room).Error
	if err != nil {
		return nil, err
	}
	return room, nil
}

func ListRooms(db *gorm.DB) (*[]model.Room, error) {
	rooms := new([]model.Room)
	err := db.Find(&rooms).Error
	if err != nil {
		return nil, err
	}
	return rooms, nil
}

func FindRoomById(db *gorm.DB, id uint) (*model.Room, error) {
	room := new(model.Room)
	err := db.First(room, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return room, nil
}
