package services

import (
	"chatroom/internal/pkg/broker"
	integration "chatroom/internal/pkg/integration/stocks"
	"encoding/json"
)

type StockService struct {
	stocksIntegration *integration.StocksIntegration
	broker            *broker.Broker
}

func NewStockService(broker *broker.Broker) *StockService {
	return &StockService{
		stocksIntegration: &integration.StocksIntegration{},
	}
}

type GetStockPriceRequest struct {
	RoomId    uint   `json:"room_id"`
	StockCode string `json:"stock_code"`
}

func (ss *StockService) ConsumeGetStockPriceRequest(msg []byte) (float64, error) {
	request := &GetStockPriceRequest{}
	err := json.Unmarshal(msg, request)
	if err != nil {
		return 0, err
	}
	value, err := ss.getStockPrice(request.StockCode)
	if err != nil {
		return 0, err
	}
	return value, nil
}

func (ss *StockService) getStockPrice(stockCode string) (float64, error) {
	value, err := ss.stocksIntegration.GetStockPrice(stockCode)
	if err != nil {
		return 0, err
	}
	return value, err
}
