package nasdaq

import (
	"fmt"
	"testing"
	"time"

	"github.com/CaledoniaProject/gopkg-commons/commons"
	"github.com/sirupsen/logrus"
)

func init() {
	fmt.Println("test init() called")

	logrus.SetLevel(logrus.DebugLevel)
	commons.SetGlobalHTTPProxy("http://127.0.0.1:6000")
}

func TestChart(t *testing.T) {
	data, err := GetHistoricalChart("AAPL", time.Now().AddDate(-20, 0, -15), time.Now().AddDate(-20, 0, 0))
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}

	for _, row := range data {
		fmt.Println(row)
	}
}

func TestQuote(t *testing.T) {
	data, err := GetQuoteChart("TLT", "etf", time.Now().AddDate(0, 0, -7), time.Now())
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}

	for _, row := range data {
		fmt.Println(row.ToHistoricalPrice())
	}
}

func TestGetStocks(t *testing.T) {
	data, err := GetStocks(5, 5)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}

	for _, row := range data {
		fmt.Println(row)
	}
}
