package chat

import (
	"chatroom/config"
	"chatroom/internal/apps/chat/consumer"
	"chatroom/internal/apps/chat/router"
	"chatroom/internal/models"
	"chatroom/internal/pkg/cache"
	"chatroom/internal/pkg/database"
	"chatroom/internal/pkg/exchanger"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
	"github.com/streadway/amqp"
	"gorm.io/gorm"
	"log/slog"
	"os"
)

type App struct {
	Router    *mux.Router
	DB        *gorm.DB
	Cache     *redis.Client
	Exchanger *amqp.Connection
	Consumer  *consumer.Consumer
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

	a.Cache = cache.NewCache(config.GetConfig())
	a.DB = database.NewDB(config.GetConfig())
	a.migrate()
	a.Router = router.NewRouter(a.DB, a.Cache).Router
	a.Exchanger = exchanger.NewExchanger(config.GetConfig())
	a.Consumer = consumer.NewConsumer(a.Exchanger)
}

func (a *App) migrate() {
	err := a.DB.AutoMigrate(&models.Room{}, &models.User{}, &models.Message{})
	if err != nil {
		panic(err.Error())
	}
}
