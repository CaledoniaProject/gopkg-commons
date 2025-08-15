package eastmoney_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/CaledoniaProject/gopkg-commons/commons"
	"github.com/CaledoniaProject/gopkg-commons/commons/finance/eastmoney"
	"github.com/sirupsen/logrus"
)

func init() {
	fmt.Println("test init() called")

	logrus.SetLevel(logrus.DebugLevel)
	commons.SetGlobalHTTPProxy("http://127.0.0.1:6000")
}

func TestFundNav(t *testing.T) {
	if data, err := eastmoney.GetFundHistoryNav(&eastmoney.LsjzRequest{
		FundCode:  "000290",
		PageIndex: 1,
		PageSize:  1,
		StartDate: time.Now().Add(-7 * 24 * time.Hour),
		EndDate:   time.Now(),
	}); err != nil {
		t.Fatalf("get failed: %v", err)
	} else if len(data) == 0 {
		t.Fatalf("no data")
	} else {
		fmt.Println(data[0])
	}
}

func TestStockPrice_CN(t *testing.T) {
	quoteTable, err := eastmoney.GetQuoteByRawCode("513400", "")
	if err != nil {
		t.Fatalf("get quote: %v", err)
	}

	prices, err := eastmoney.GetHistoryPrice(quoteTable.QuoteID, time.Now().Add(-3*24*time.Hour), time.Now())
	if err != nil {
		t.Fatalf("get price: %v", err)
	}

	for _, price := range prices {
		fmt.Println(price)
	}
}

func TestStockPrice_HK(t *testing.T) {
	quoteTable, err := eastmoney.GetQuoteByRawCode("00700", "")
	if err != nil {
		t.Fatalf("get quote: %v", err)
	}

	prices, err := eastmoney.GetHistoryPrice(quoteTable.QuoteID, time.Now().Add(-3*24*time.Hour), time.Now())
	if err != nil {
		t.Fatalf("get price: %v", err)
	}

	for _, price := range prices {
		fmt.Println(price)
	}
}

func TestStockPrice_US(t *testing.T) {
	quoteTable, err := eastmoney.GetQuoteByRawCode("BIL", "")
	if err != nil {
		t.Fatalf("get quote: %v", err)
	}

	prices, err := eastmoney.GetHistoryPrice(quoteTable.QuoteID, time.Now().Add(-3*24*time.Hour), time.Now())
	if err != nil {
		t.Fatalf("get price: %v", err)
	}

	for _, price := range prices {
		fmt.Println(price)
	}
}
