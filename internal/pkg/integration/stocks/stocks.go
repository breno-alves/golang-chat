package integration

import (
	"encoding/csv"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strconv"
)

type StocksIntegration struct{}

var (
	URL = "https://stooq.com"
)

func (si *StocksIntegration) GetStockPrice(stockCode string) (float64, error) {
	url := fmt.Sprintf("%s/q/l/?s=%s&f=sd2t2ohlcv&h&e=csv", URL, stockCode)
	response, err := http.Get(url)
	if err != nil {
		return 0, err
	}
	r := csv.NewReader(response.Body)
	header, err := r.Read()
	if err != nil {
		return 0, err
	}
	value, err := si.processCSV(r, header)
	if err != nil {
		return 0, err
	}
	return value, nil
}

func (si *StocksIntegration) processCSV(reader *csv.Reader, _ []string) (float64, error) {
	var closeValue float64
	record, err := reader.Read()
	if err == io.EOF {
	} else if err != nil {
		slog.Error(fmt.Sprintf("Error reading CSV data: %s", err))
	}
	closeValue, err = strconv.ParseFloat(record[6], 8)
	if err != nil {
		slog.Error(fmt.Sprintf("Error parsing CSV data: %s", err))
		return 0, err
	}
	return closeValue, nil
}
