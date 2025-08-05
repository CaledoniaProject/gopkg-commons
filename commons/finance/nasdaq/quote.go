package nasdaq

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/CaledoniaProject/gopkg-commons/commons"
	"github.com/CaledoniaProject/gopkg-commons/commons/finance"
)

type QuoteResponse struct {
	Status *QuoteStatus `json:"status"`
	Data   *QuoteData   `json:"data"`
}

type QuoteStatus struct {
	RCode int `json:"rCode"`
}

type QuoteData struct {
	Chart []*QuoteDataChart `json:"chart"`
}

type QuoteDataChart struct {
	X int          `json:"x"`
	Y float64      `json:"y"`
	Z *QuoteDetail `json:"z"`
}

type QuoteDetail struct {
	High     string `json:"high"`
	Low      string `json:"low"`
	Open     string `json:"open"`
	Close    string `json:"close"`
	Volume   string `json:"volume"`
	DateTime string `json:"dateTime"`
	Value    string `json:"value"`
}

func (q *QuoteDetail) RemoveComma() {
	q.High = strings.Replace(q.High, ",", "", -1)
	q.Low = strings.Replace(q.Low, ",", "", -1)
	q.Open = strings.Replace(q.Open, ",", "", -1)
	q.Close = strings.Replace(q.Close, ",", "", -1)
	q.Volume = strings.Replace(q.Volume, ",", "", -1)
	q.Value = strings.Replace(q.Value, ",", "", -1)
}

func (q *QuoteDetail) ToHistoricalPrice() (*finance.HistoricalPrice, error) {
	var (
		result = &finance.HistoricalPrice{}
	)

	// high
	if tmp, err := strconv.ParseFloat(q.High, 64); err != nil {
		return nil, err
	} else {
		result.HighPrice = tmp
	}

	// low
	if tmp, err := strconv.ParseFloat(q.Low, 64); err != nil {
		return nil, err
	} else {
		result.LowPrice = tmp
	}

	// open
	if tmp, err := strconv.ParseFloat(q.Open, 64); err != nil {
		return nil, err
	} else {
		result.OpenPrice = tmp
	}

	// close
	if tmp, err := strconv.ParseFloat(q.Close, 64); err != nil {
		return nil, err
	} else {
		result.ClosePrice = tmp
	}

	// date
	if tmp, err := time.Parse("1/2/2006", q.DateTime); err != nil {
		return nil, err
	} else {
		result.Date = tmp
	}

	return result, nil
}

func GetQuoteChart(symbol, assetClass string, fromDate, toDate time.Time) (quotes []*QuoteDetail, err error) {
	var (
		quoteResponse = &QuoteResponse{}
		quoteURL      = fmt.Sprintf("https://api.nasdaq.com/api/quote/%s/chart?assetclass=%s&fromdate=%s&todate=%s",
			url.QueryEscape(symbol), url.QueryEscape(assetClass), fromDate.Format("2006-01-02"), toDate.Format("2006-01-02"),
		)
	)

	_, body, err := commons.HttpRequest(&commons.RequestOptions{
		URL:         quoteURL,
		MaxBodyRead: 10 * 1024 * 1024,
		Timeout:     30,
		Headers: map[string]string{
			"User-Agent": commons.RandomUserAgent(),
		},
	})
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(body, quoteResponse); err != nil {
		return nil, err
	} else if quoteResponse.Data == nil || len(quoteResponse.Data.Chart) == 0 {
		if quoteResponse.Status != nil {
			return nil, fmt.Errorf("no result, rCode=%d", quoteResponse.Status.RCode)
		} else {
			return nil, errors.New("no result")
		}
	}

	for _, chart := range quoteResponse.Data.Chart {
		chart.Z.RemoveComma()
		quotes = append(quotes, chart.Z)
	}

	return quotes, nil
}
