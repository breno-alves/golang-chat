package services

import (
	"chatroom/internal/models"
	"gorm.io/gorm"
)

func CreateRoom(db *gorm.DB) (*models.Room, error) {
	room := models.NewRoom()
	err := db.Create(room).Error
	if err != nil {
		return nil, err
	}
	return room, nil
}

func ListRooms(db *gorm.DB) (*[]models.Room, error) {
	rooms := new([]models.Room)
	err := db.Find(&rooms).Error
	if err != nil {
		return nil, err
	}
	return rooms, nil
}

func FindRoomById(db *gorm.DB, id uint) (*models.Room, error) {
	room := new(models.Room)
	err := db.First(room, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return room, nil
}
