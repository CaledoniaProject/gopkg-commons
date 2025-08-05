package eastmoney

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/CaledoniaProject/gopkg-commons/commons"
)

var (
	MarketNumber2Country = map[string]string{
		"0":   "cn",
		"1":   "cn",
		"116": "hk",
		"105": "us",
		"106": "us",
		"107": "us",
		"153": "us",
	}
)

type StockSuggestRequest struct {
	Input    string `param:"input"`
	Type     int    `param:"type"`
	Count    int    `param:"count"`
	Classify string `param:"classify"`
}

type EastMoneyQuotationCodeTableData struct {
	Code             string `json:"Code" gorm:"size:50"`
	Name             string `json:"Name" gorm:"size:512"`
	PinYin           string `json:"PinYin" gorm:"size:512"`
	ID               string `json:"ID" gorm:"size:256"`
	JYS              string `json:"JYS" gorm:"size:256"`
	Classify         string `json:"Classify" gorm:"size:32"` // 可选 AStock/HK/UsStock/OTCBB/23
	MarketType       string `json:"MarketType" gorm:"size:10"`
	SecurityTypeName string `json:"SecurityTypeName" gorm:"size:20"`
	SecurityType     string `json:"SecurityType" gorm:"size:10"`
	MktNum           string `json:"MktNum" gorm:"size:10"`
	TypeUS           string `json:"TypeUS" gorm:"size:10"`
	QuoteID          string `json:"QuoteID" gorm:"primaryKey"`
	UnifiedCode      string `json:"UnifiedCode" gorm:"size:32"`
	InnerCode        string `json:"InnerCode" gorm:"size:32"`
}

type EastMoneyQuotationCodeTableResult struct {
	Data    []*EastMoneyQuotationCodeTableData `json:"Data"`
	Status  int                                `json:"Status"`
	Message string                             `json:"Message"`
}

type StockSuggestResponse struct {
	QuotationCodeTable *EastMoneyQuotationCodeTableResult `json:"QuotationCodeTable"`
}

func GetSuggestedStock(suggestRequest *StockSuggestRequest) (result []*EastMoneyQuotationCodeTableData, err error) {
	var (
		suggestURL = fmt.Sprintf("https://searchadapter.eastmoney.com/api/suggest/get?%s", commons.EncodeStructToURLParams(suggestRequest))
		response   = &StockSuggestResponse{}
	)

	_, body, err := commons.HttpRequest(&commons.RequestOptions{
		URL:         suggestURL,
		ReadBody:    true,
		MaxBodyRead: 10 * 1024 * 1024,
		Timeout:     10,
		MaxRetry:    5,
		Headers: map[string]string{
			"Conection":  "close",
			"User-Agent": commons.RandomUserAgent(),
		},
	})

	if err != nil {
		return nil, err
	} else if err = json.Unmarshal(body, response); err != nil {
		return nil, err
	} else if response.QuotationCodeTable == nil {
		return nil, errors.New("bad json response, no data returned")
	} else if response.QuotationCodeTable.Status != 0 {
		return nil, fmt.Errorf("eastmoney api failed: %s (%d)", response.QuotationCodeTable.Message, response.QuotationCodeTable.Status)
	}

	return response.QuotationCodeTable.Data, nil
}

func GetQuoteByRawCode(rawCode, classify string) (result *EastMoneyQuotationCodeTableData, err error) {
	candidates, err := GetSuggestedStock(&StockSuggestRequest{
		Input:    rawCode,
		Classify: classify,
		Type:     14,
		Count:    100,
	})
	if err != nil {
		return nil, err
	}

	for _, row := range candidates {
		if row.Code == rawCode {
			return row, nil
		}
	}

	return nil, errors.New("no matched code found")
}
