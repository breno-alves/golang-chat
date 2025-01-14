package router

import (
	"chatroom/internal/handlers/bot"
	"chatroom/internal/pkg/broker"
)

var (
	BotRequestQueue  = "BOT_STOCKS"
	BotResponseQueue = "BOT_STOCKS_RESPONSE"
)

type Router struct {
	handler *bot.Handler
}

func NewRouter(broker *broker.Broker) *Router {
	r := &Router{
		handler: bot.NewHandler(),
	}
	go r.InitRmqRouter(broker)
	return r
}
