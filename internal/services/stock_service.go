package services

import integration "chatroom/internal/pkg/integration/stocks"

type StockService struct {
	stocksIntegration *integration.StocksIntegration
}

func NewStockService() *StockService {
	return &StockService{
		stocksIntegration: &integration.StocksIntegration{},
	}
}

func (ss *StockService) GetStockPrice(stockCode string) (float64, error) {
	value, err := ss.stocksIntegration.GetStockPrice(stockCode)
	if err != nil {
		return 0, err
	}
	return value, err
}
