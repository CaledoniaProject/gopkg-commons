package eastmoney

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"sort"
	"strings"
	"time"

	"github.com/CaledoniaProject/gopkg-commons/commons"
)

// 不确定的: f85 也是总股本，f117 也是总市值
type qtStockGetResponseData struct {
	Price             float64 `json:"f43"`
	OutstandingShares float64 `json:"f84"` // 总股本
	NAV_TTM           float64 `json:"f92"`
	EPS_TTM           float64 `json:"f108"`
	MarketCap         float64 `json:"f116"` // 总市值
	DividendYield     float64 `json:"f126"`
	PE_TTM            float64 `json:"f164"`
	PS_TTM            float64 `json:"f165"`
	PB_TTM            float64 `json:"f167"`
	FiftyTwoHigh      float64 `json:"f174"`
	FiftyTwoLow       float64 `json:"f175"`
}

type qtStockGetResponse struct {
	Data *qtStockGetResponseData `json:"data"`
}

// 获取股票基本信息
func GetSnapshot(quoteId string) (*qtStockGetResponseData, error) {
	var (
		pushURL = fmt.Sprintf("https://push2.eastmoney.com/api/qt/stock/get?secid=%s&fields=f165,f126,f174,f175,f58,f107,f57,f43,f59,f169,f170,f152,f46,f60,f44,f45,f47,f48,f19,f532,f39,f161,f49,f171,f50,f86,f600,f601,f154,f84,f85,f168,f108,f116,f167,f164,f92,f71,f117,f292,f301", url.QueryEscape(quoteId))
		result  = &qtStockGetResponse{}
	)

	_, body, err := commons.HttpRequest(&commons.RequestOptions{
		URL:         pushURL,
		ReadBody:    true,
		MaxBodyRead: 10 * 1024 * 1024,
		Timeout:     10,
		MaxRetry:    5,
		MaxRedirect: 1,
		Headers: map[string]string{
			"User-Agent": commons.RandomUserAgent(),
		},
	})
	if err != nil {
		return nil, err
	} else if err = json.Unmarshal(body, result); err != nil {
		return nil, err
	} else if result.Data == nil {
		return nil, fmt.Errorf("no metadata available for %s", quoteId)
	}

	return result.Data, nil
}

type qtStockKLineGetRequest struct {
	Secid     string   `param:"secid"`
	Fields1   []string `param:"fields1"`
	Fields2   []string `param:"fields2"`
	KLT       int      `param:"klt"`
	FQT       int      `param:"fqt"`
	BeginDate string   `param:"beg"`
	EndDate   string   `param:"end"`
}

type qtStockKLineGetResponseData struct {
	Code   string   `json:"code"`   // 原始代码，e.g 600000
	Market int      `json:"market"` // 市场编号
	KLines []string `json:"klines"` // 格式类似 2023-12-08 09:31,6.60
}

type qtStockKLineGetResponse struct {
	Data *qtStockKLineGetResponseData `json:"data"`
}

type historyPrice struct {
	Date        time.Time
	OpenPrice   float64
	ClosePrice  float64
	HighPrice   float64
	LowPrice    float64
	TradeVolume float64
	TradeAmount float64
}

// f51: 日期
// f52: 开盘价
// f53: 收盘价
// f54: 最高价
// f55: 最低价
// f56: 交易量
// f57: 成交额
// 日期格式: 20990101
func GetHistoryPrice(quoteId string, beginDate, endDate time.Time) (result []*historyPrice, err error) {
	var (
		historyPriceParams = &qtStockKLineGetRequest{
			Secid:     quoteId,
			Fields1:   []string{"f1", "f2", "f3", "f4", "f5", "f6"},
			Fields2:   []string{"f51", "f52", "f53", "f54", "f55", "f56", "f57"},
			KLT:       101,
			FQT:       1,
			BeginDate: beginDate.Format("20060102"),
			EndDate:   endDate.Format("20060102"),
		}
		historyPriceURL = fmt.Sprintf("https://push2his.eastmoney.com/api/qt/stock/kline/get?%s", commons.EncodeStructToURLParams(historyPriceParams))
		historyResponse = &qtStockKLineGetResponse{}
	)

	_, body, err := commons.HttpRequest(&commons.RequestOptions{
		URL:         historyPriceURL,
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
	} else if err = json.Unmarshal(body, historyResponse); err != nil {
		return nil, err
	} else if historyResponse.Data == nil {
		return nil, errors.New("no data available")
	}

	for _, row := range historyResponse.Data.KLines {
		parts := strings.Split(row, ",")
		if len(parts) != 7 {
			continue
		}

		tradeDay, err := time.Parse("2006-01-02", parts[0])
		if err != nil {
			return nil, err
		}

		result = append(result, &historyPrice{
			Date:        tradeDay,
			OpenPrice:   commons.Atof(parts[1]),
			ClosePrice:  commons.Atof(parts[2]),
			HighPrice:   commons.Atof(parts[3]),
			LowPrice:    commons.Atof(parts[4]),
			TradeVolume: commons.Atof(parts[5]),
			TradeAmount: commons.Atof(parts[6]),
		})
	}

	// 日期从小到大
	sort.Slice(result, func(i, j int) bool {
		return result[i].Date.Before(result[j].Date)
	})
	return result, nil
}
