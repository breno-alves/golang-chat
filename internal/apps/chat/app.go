package chat

import (
	"chatroom/config"
	"chatroom/internal/apps/chat/router"
	"chatroom/internal/models"
	"chatroom/internal/pkg/broker"
	"chatroom/internal/pkg/cache"
	"chatroom/internal/pkg/database"
	"errors"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"log/slog"
	"os"
)

type App struct {
	Router *mux.Router
	DB     *gorm.DB
	Cache  *redis.Client
	Broker *broker.Broker
}

func NewApp() *App {
	app := &App{}
	app.Initialize()
	return app
}

func (a *App) Initialize() {
	logger := slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelDebug})
	slog.SetDefault(slog.New(logger))

	err := godotenv.Load(".env")
	if err != nil {
		slog.Debug("app did not load .env")
	}

	a.Broker = broker.NewBroker(config.GetConfig())
	a.Cache = cache.NewCache(config.GetConfig())
	a.DB = database.NewDB(config.GetConfig())
	a.migrate()

	a.Router = router.NewRouter(a.DB, a.Cache, a.Broker).HttpRouter

}

func (a *App) migrate() {
	err := a.DB.AutoMigrate(&models.Room{}, &models.Message{}, &models.User{})
	if err != nil {
		panic(err.Error())
	}

	// Creates bot user if it doesn't exist
	if errors.Is(a.DB.First(&models.User{}, "username = ?", "bot").Error, gorm.ErrRecordNotFound) {
		a.DB.Create(&models.User{Username: "bot", Password: "123456"})
	}

	// Creates first room if it doesn't exist
	if errors.Is(a.DB.First(&models.Room{}).Error, gorm.ErrRecordNotFound) {
		a.DB.Create(&models.Room{Title: "First room"})
	}
}
