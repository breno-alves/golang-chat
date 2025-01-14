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

func NewStockService() *StockService {
	return &StockService{
		stocksIntegration: &integration.StocksIntegration{},
	}
}

type GetStockPriceRequest struct {
	RoomId    uint   `json:"room_id"`
	StockCode string `json:"stock_code"`
}

type GetStockPriceResponse struct {
	RoomId    uint    `json:"room_id"`
	StockCode string  `json:"stock_code"`
	Value     float64 `json:"value"`
}

func (ss *StockService) ConsumeGetStockPrice(msg []byte) ([]byte, error) {
	request := &GetStockPriceRequest{}
	err := json.Unmarshal(msg, request)
	if err != nil {
		return nil, err
	}
	value, err := ss.getStockPrice(request.StockCode)
	if err != nil {
		return nil, err
	}
	response, err := json.Marshal(&GetStockPriceResponse{
		RoomId:    request.RoomId,
		StockCode: request.StockCode,
		Value:     value,
	})
	if err != nil {
		return nil, err
	}
	return response, nil
}

func (ss *StockService) getStockPrice(stockCode string) (float64, error) {
	value, err := ss.stocksIntegration.GetStockPrice(stockCode)
	if err != nil {
		return 0, err
	}
	return value, err
}
