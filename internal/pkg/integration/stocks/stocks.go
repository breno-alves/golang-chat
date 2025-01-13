package integration

import (
	"encoding/csv"
	"fmt"
	"io"
	"net/http"
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
	si.processCSV(r, header)

	return 0, nil
}

func (si *StocksIntegration) processCSV(reader *csv.Reader, header []string) {
	fmt.Println(header)
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			fmt.Println("Error reading CSV data:", err)
			break
		}
		for i, field := range header {
			fmt.Printf("%s: %s\n", field, record[i])
		}
	}
}
