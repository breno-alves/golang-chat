package router

import (
	"chatroom/internal/pkg/broker"
	"fmt"
	"log/slog"
)

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
	go r.GetStockPriceResponse(broker)
}

func (r *Router) GetStockPriceResponse(broker *broker.Broker) {
	// This will receive response from the bot API
	channel, err := broker.NewConsumer(BotResponseQueue)
	if err != nil {
		panic(err)
	}
	for {
		select {
		case msg := <-channel.Ch:
			slog.Debug(fmt.Sprintf("received message: %s", string(msg)))
			message, err := r.handler.MessageService.CreateMessageFromBot(msg)
			if err != nil {
				slog.Error("could not create a message")
				continue
			}
			err = r.handler.BroadcastMessage(nil, message)
			if err != nil {
				slog.Error("could not broadcast a message")
				continue
			}
		}
	}
}
