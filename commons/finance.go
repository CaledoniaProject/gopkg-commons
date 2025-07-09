package commons

import (
	"fmt"
	"math"
)

// 计算复合收益率
func CompoundInterestRate(startValue, endValue float64, years int) float64 {
	if years == 0 || startValue <= 0 || endValue <= 0 {
		return 0
	}

	rate := math.Pow(endValue/startValue, 1.0/float64(years)) - 1
	return rate
}

// 复合增长率不低于目标
func IsCAGRBiggerThan(input []float64, minCAGR float64) bool {
	if len(input) < 2 {
		return false
	}

	cagr := math.Pow(input[0]/input[len(input)-1], 1.0/float64(len(input)))
	return cagr > minCAGR
}

// 年度增长大于目标
func IsAllYoYBiggerThan(input []float64, yoy float64) bool {
	for i := 0; i < len(input)-1; i++ {
		rate := (input[i] - input[i+1]) / input[i+1]
		if rate < yoy {
			return false
		}
	}

	return len(input) > 0
}

// 货币格式化
func FormatCurrency(val float64) string {
	var (
		negative = false
		val2     = val
		suffix   = ""
	)

	if val2 < 0 {
		negative = true
		val2 = -val2
	}

	if val2 > math.Pow10(9) {
		val2 = val2 / math.Pow10(9)
		suffix = " B"
	} else if val2 > math.Pow10(6) {
		val2 = val2 / math.Pow10(6)
		suffix = " M"
	}

	if negative {
		val2 = -val2
	}

	return fmt.Sprintf("%.2f%s", val2, suffix)
}
