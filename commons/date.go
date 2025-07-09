package commons

import (
	"time"
)

// 获取第一天
func TruncateToStartOfDay(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
}

// 获取最后一天最后一秒
func TruncateToEndOfDay(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 23, 59, 59, 0, t.Location())
}

// 上个季度最后1天
func GetLastQuarterEndDate() time.Time {
	var (
		result = time.Now()
		month  = result.Month()
	)

	if month >= 1 && month <= 3 {
		return time.Date(result.Year()-1, 12, 31, 0, 0, 0, 0, result.Location())
	} else if month >= 4 && month <= 6 {
		return time.Date(result.Year(), 3, 31, 0, 0, 0, 0, result.Location())
	} else if month >= 7 && month <= 9 {
		return time.Date(result.Year(), 6, 30, 0, 0, 0, 0, result.Location())
	}

	return time.Date(result.Year(), 9, 30, 0, 0, 0, 0, result.Location())
}

// 这个季度最后1天
func GetCurrentQuarterEndDate() time.Time {
	var (
		result = time.Now()
		month  = result.Month()
	)

	if month >= 1 && month <= 3 {
		return time.Date(result.Year(), 3, 31, 0, 0, 0, 0, result.Location())
	} else if month >= 4 && month <= 6 {
		return time.Date(result.Year(), 6, 30, 0, 0, 0, 0, result.Location())
	} else if month >= 7 && month <= 9 {
		return time.Date(result.Year(), 9, 30, 0, 0, 0, 0, result.Location())
	}

	return time.Date(result.Year(), 12, 31, 0, 0, 0, 0, result.Location())
}

// 两个日期之间的天数
func DaysBetween(date1, date2 string) (float64, error) {
	dateObj1, err := time.Parse("2006-01-02 00:00:00", date1)
	if err != nil {
		return 0, err
	}

	dateObj2, err := time.Parse("2006-01-02 00:00:00", date2)
	if err != nil {
		return 0, err
	}

	return dateObj1.Sub(dateObj2).Hours() / 24, nil
}

// 计算天数
func DaysToNow(dateStr string) (float64, error) {
	dateObj, err := time.Parse("2006-01-02 00:00:00", dateStr)
	if err != nil {
		return 0, err
	}

	return time.Since(dateObj).Hours() / 24, nil
}
