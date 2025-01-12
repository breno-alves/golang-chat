package database

import (
	"chatroom/config"
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log/slog"
)

func NewDB(config *config.Config) *gorm.DB {
	defer func() {
		slog.Info("successfully connected to database")
	}()

	dns := fmt.Sprintf("host=%s port=%d user=%s dbname=%s sslmode=disable password=%s", config.DB.Host, config.DB.Port, config.DB.Username, config.DB.Name, config.DB.Password)
	slog.Debug(fmt.Sprintf("connecting to database: %s", dns))

	db, err := gorm.Open(postgres.Open(dns), &gorm.Config{})
	if err != nil {
		panic("could not connect to database")
		return nil
	}

	return db
}
