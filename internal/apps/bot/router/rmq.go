package router

import (
	"chatroom/internal/pkg/broker"
	"fmt"
	"log/slog"
)

func (r *Router) InitRmqRouter(broker *broker.Broker) {
	// This will send response to chat app
	_, err := broker.NewProducer(BotResponseQueue)
	if err != nil {
		panic(err)
	}
	go r.GetStockPriceRequest(broker)
}

func (r *Router) GetStockPriceRequest(broker *broker.Broker) {
	// This will consume request from to chat app
	channel, err := broker.NewConsumer(BotRequestQueue)
	if err != nil {
		panic(err)
	}
	for {
		select {
		case msg := <-*channel.Ch:
			slog.Debug(fmt.Sprintf("Received a message: %s", string(msg)))
		}
	}
}
