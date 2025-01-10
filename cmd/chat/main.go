package main

import (
	"chatroom/internal/chat/cache"
	"chatroom/internal/chat/config"
	"chatroom/internal/chat/database"
	"chatroom/internal/chat/router"
	"github.com/gorilla/mux"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"log"
	"log/slog"
	"net/http"
	"os"
)

type App struct {
	Router *mux.Router
	DB     *gorm.DB
	Cache  *redis.Client
}

func (a *App) Initialize() {
	a.DB = database.NewDB(config.GetConfig())
	a.Cache = cache.NewCache(config.GetConfig())
	a.Router = router.NewRouter(a.DB, a.Cache).Router
}

func (a *App) Run() {
	err := http.ListenAndServe(":8080", a.Router)
	if err != nil {
		log.Fatal(err)
	}
}
func main() {
	logger := slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelDebug})
	slog.SetDefault(slog.New(logger))
	app := &App{}
	app.Initialize()
	app.Run()
}
