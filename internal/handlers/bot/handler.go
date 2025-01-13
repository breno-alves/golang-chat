package bot

import (
	"chatroom/internal/pkg/broker"
	"chatroom/internal/services"
)

type Handler struct {
	stockService *services.StockService
}

func NewHandler(broker *broker.Broker) *Handler {
	return &Handler{
		stockService: services.NewStockService(),
	}
}
