package nasdaq

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"time"

	"github.com/CaledoniaProject/gopkg-commons/commons"
)

type chartResponse struct {
	CompanyName string        `json:"companyName"`
	MarketData  []*marketData `json:"marketData"`
}

type marketData struct {
	High   float64 `json:"high"`
	Low    float64 `json:"low"`
	Open   float64 `json:"open"`
	Close  float64 `json:"close"`
	Volume float64 `json:"volume"`
	Date   string  `json:"date"`
}

func GetHistoricalChart(symbol string, startDate, endDate time.Time) (results []*marketData, err error) {
	var (
		chartResponse = &chartResponse{}
		chartURL      = fmt.Sprintf("https://charting.nasdaq.com/data/charting/historical?symbol=%s&date=%s~%s",
			url.QueryEscape(symbol),
			startDate.Format("2006-01-02"),
			endDate.Format("2006-01-02"),
		)
	)

	_, body, err := commons.HttpRequest(&commons.RequestOptions{
		URL:         chartURL,
		MaxBodyRead: 10 * 1024 * 1024,
		Timeout:     30,
		Headers: map[string]string{
			"User-Agent": commons.RandomUserAgent(),
			"Referer":    "https://charting.nasdaq.com/dynamic/chart.html",
		},
	})
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(body, chartResponse); err != nil {
		return nil, err
	} else if len(chartResponse.MarketData) == 0 {
		return nil, errors.New("no result")
	}

	return chartResponse.MarketData, nil
}
