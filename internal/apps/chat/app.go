package chat

import (
	"chatroom/config"
	"chatroom/internal/apps/chat/router"
	"chatroom/internal/models"
	"chatroom/internal/pkg/cache"
	"chatroom/internal/pkg/database"
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

	a.DB = database.NewDB(config.GetConfig())
	a.migrate()

	a.Cache = cache.NewCache(config.GetConfig())
	a.Router = router.NewRouter(a.DB, a.Cache).Router
}

func (a *App) migrate() {
	err := a.DB.AutoMigrate(&models.Room{}, &models.User{}, &models.Message{})
	if err != nil {
		panic(err.Error())
	}
}
