package bot

import (
	"chatroom/internal/services"
)

type Handler struct {
	StockService *services.StockService
}

func NewHandler() *Handler {
	return &Handler{
		StockService: services.NewStockService(),
	}
}
