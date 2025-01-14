package router

import (
	"chatroom/internal/pkg/broker"
	"errors"
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
			response, err := r.handler.StockService.ConsumeGetStockPrice(msg)
			if err != nil {
				slog.Error(fmt.Sprintf("error consuming message: %s", err.Error()))
				continue
			}

			err = r.SendStockPriceResponse(broker, response)
			if err != nil {
				slog.Error(fmt.Sprintf("error sending response: %s", err.Error()))
				continue
			}
		}
	}
}

func (r *Router) SendStockPriceResponse(broker *broker.Broker, response []byte) error {
	channel, ok := broker.Channels[BotResponseQueue]
	if !ok {
		slog.Error("could not find broker channel")
		return errors.New("could not find broker channel")
	}
	channel.Ch <- response
	return nil
}
