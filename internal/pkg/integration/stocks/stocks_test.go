package integration

import (
	"encoding/csv"
	"fmt"
	"strings"
	"testing"
)

func TestProcessCSV(t *testing.T) {
	input := "AAPL.US,2025-01-10,22:00:16,240.01,240.16,233,236.85,61710856\n"

	si := &StocksIntegration{}
	reader := csv.NewReader(strings.NewReader(input))

	value, err := si.processCSV(reader, []string{})
	if err != nil {
		t.Errorf("error in process csv %s", err.Error())
	}
	fmt.Println(value)
}
