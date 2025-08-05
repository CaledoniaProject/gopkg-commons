package finance

import "time"

type HistoricalPrice struct {
	Date        time.Time
	OpenPrice   float64
	ClosePrice  float64
	HighPrice   float64
	LowPrice    float64
	TradeVolume float64
	TradeAmount float64
}
