package bot

import "chatroom/internal/services"

type Handler struct {
	stockService *services.StockService
}

func NewHandler() *Handler {
	return &Handler{
		stockService: services.NewStockService(),
	}
}
