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
		case msg := <-channel.Ch:
			fmt.Println("received bot:", string(msg))
			value, err := r.handler.StockService.ConsumeGetStockPriceRequest(msg)
			if err != nil {
				slog.Error(err.Error())
			}
			slog.Debug(fmt.Sprintf("stock value is: %v", value))
			continue
		default:
			continue
		}
	}
}
