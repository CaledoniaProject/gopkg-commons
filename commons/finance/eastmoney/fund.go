package eastmoney

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/CaledoniaProject/gopkg-commons/commons"
)

type LsjzRequest struct {
	FundCode  string    `param:"fundCode"`
	PageIndex int       `param:"pageIndex"`
	PageSize  int       `param:"pageSize"`
	StartDate time.Time `param:"startDate" layout:"2006-01-02"`
	EndDate   time.Time `param:"endDate" layout:"2006-01-02"`
}

type LsjzResponse struct {
	Data       LsjzResponseData `json:"Data"`
	ErrCode    int              `json:"ErrCode"`
	ErrMsg     string           `json:"ErrMsg"`
	TotalCount int              `json:"TotalCount"`
}

type LsjzResponseData struct {
	FundNavList []*FundNav `json:"LSJZList"`
}

type FundNav struct {
	FSRQ      string `json:"FSRQ"`      // 净值日期
	DWJZ      string `json:"DWJZ"`      // 单位净值
	LJJZ      string `json:"LJJZ"`      // 累计净值
	JZZZL     string `json:"JZZZL"`     // 日增长率
	ACTUALSYI string `json:"ACTUALSYI"` // 实际收益
	SGZT      string `json:"SGZT"`      // 申购状态
	SHZT      string `json:"SHZT"`      // 赎回状态
	NAVTYPE   string `json:"NAVTYPE"`   // 净值类型
	FHFCZ     string `json:"FHFCZ"`     // 分红发放差值
	FHFCBZ    string `json:"FHFCBZ"`    // 分红发放标准
	FHSP      string `json:"FHSP"`      // 分红说明
	SDATE     string `json:"SDATE"`     // 未知
	DTYPE     string `json:"DTYPE"`     // 未知

	UnitValue float64   // 手动计算的净值
	UnitDate  time.Time // 解析后的日期
}

func GetFundHistoryNav(lsjzRequest *LsjzRequest) (results []*FundNav, err error) {
	var (
		lsjzResponse = &LsjzResponse{}
	)

	if _, body, err := commons.HttpRequest(&commons.RequestOptions{
		URL: "https://api.fund.eastmoney.com/f10/lsjz?" + commons.EncodeStructToURLParams(lsjzRequest),
		Headers: map[string]string{
			"User-Agent": commons.RandomUserAgent(),
			"Referer":    "https://fundf10.eastmoney.com/",
		},
		MaxBodyRead: 10 * 1024 * 1024,
	}); err != nil {
		return nil, err
	} else if err := json.Unmarshal(body, lsjzResponse); err != nil {
		return nil, err
	} else if lsjzResponse.ErrCode != 0 {
		return nil, fmt.Errorf("error getting nav for %s: %v (%d)", lsjzRequest.FundCode, lsjzResponse.ErrMsg, lsjzResponse.ErrCode)
	} else {
		for _, row := range lsjzResponse.Data.FundNavList {
			// 忽略错误
			if tmp, err := strconv.ParseFloat(row.DWJZ, 64); err == nil {
				row.UnitValue = tmp
			}

			if tmp, err := time.Parse("2006-01-02", row.FSRQ); err == nil {
				row.UnitDate = tmp
			}
		}

		return lsjzResponse.Data.FundNavList, nil
	}
}
