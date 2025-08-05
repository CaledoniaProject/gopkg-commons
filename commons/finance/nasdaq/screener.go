package nasdaq

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/CaledoniaProject/gopkg-commons/commons"
)

type ScreenerResponse struct {
	Data *ScreenerData `json:"data"`
}

type ScreenerData struct {
	Table *ScreenerTable `json:"table"`
}

type ScreenerTable struct {
	Rows []*StockRow `json:"rows"`
}

type StockRow struct {
	Symbol        string `json:"symbol"`
	Name          string `json:"name"`
	LastScale     string `json:"lastsale"`
	NetChange     string `json:"netchange"`
	PercentChange string `json:"pctchange"`
	MarketCap     string `json:"marketCap"`
	EarningsDate  string `json:"earningsDate"`
}

func GetStocks(limit, offset int) (result []*StockRow, err error) {
	var (
		resp = &ScreenerResponse{}
	)

	_, body, err := commons.HttpRequest(&commons.RequestOptions{
		URL:         fmt.Sprintf("https://api.nasdaq.com/api/screener/stocks?tableonly=true&limit=%d&offset=%d", limit, offset),
		MaxBodyRead: 10 * 1024 * 1024,
		Timeout:     30,
		Headers: map[string]string{
			"User-Agent": commons.RandomUserAgent(),
		},
	})
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(body, resp); err != nil {
		return nil, err
	} else if resp.Data != nil && resp.Data.Table != nil {
		return resp.Data.Table.Rows, nil
	}

	return nil, errors.New("empty fields")
}
