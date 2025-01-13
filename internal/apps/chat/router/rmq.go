package router

import "chatroom/internal/pkg/broker"

var (
	BotRequestQueue  = "BOT_STOCKS"
	BotResponseQueue = "BOT_STOCKS_RESPONSE"
)

func (r *Router) InitRmqRouter(broker *broker.Broker) {
	// This will send request to bot API to retrieve stocks data
	_, err := broker.NewProducer(BotRequestQueue)
	if err != nil {
		panic(err)
	}

	// This will receive response from the bot API
	_, err = broker.NewConsumer(BotResponseQueue)
	if err != nil {
		panic(err)
	}
}
