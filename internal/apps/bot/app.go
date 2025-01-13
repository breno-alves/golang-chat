package bot

import (
	"chatroom/config"
	"chatroom/internal/apps/bot/router"
	//"chatroom/internal/apps/bot/router/exchanger/producer"
	"chatroom/internal/pkg/broker"
	//integration "chatroom/internal/pkg/integration/stocks"
	"github.com/joho/godotenv"
	"log/slog"
	"os"
)

type App struct {
	Broker *broker.Broker
}

func NewApp() *App {
	//stocksAPI := integration.StocksIntegration{}
	//err := stocksAPI.GetStock("aapl.us")
	//if err != nil {
	//	panic(err)
	//}
	//
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
	router.NewRouter(a.Broker)
	a.KeepAlive()
}

func (a *App) KeepAlive() {
	forever := make(chan bool)
	<-forever
}
