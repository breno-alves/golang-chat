package bot

import (
	"chatroom/config"
	"chatroom/internal/apps/bot/router/exchanger/producer"
	"chatroom/internal/pkg/broker"
	integration "chatroom/internal/pkg/integration/stocks"
)

type App struct {
	broker.Exchanger
}

func NewApp() *App {
	app := &App{
		Producer: producer.NewProducer(config.GetConfig()),
	}
	stocksAPI := integration.StocksIntegration{}
	err := stocksAPI.GetStock("aapl.us")
	if err != nil {
		panic(err)
	}

	return app
}
