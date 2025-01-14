package bot

import (
	"chatroom/internal/pkg/broker"
	"chatroom/internal/services"
)

type Handler struct {
	StockService *services.StockService
}

func NewHandler(broker *broker.Broker) *Handler {
	return &Handler{
		StockService: services.NewStockService(broker),
	}
}
