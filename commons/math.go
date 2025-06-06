package commons

import (
	"fmt"
	"math"
	"strconv"
)

func CompoundInterestRate(startValue, endValue float64, years int) float64 {
	if years == 0 || startValue <= 0 || endValue <= 0 {
		return 0
	}

	rate := math.Pow(endValue/startValue, 1.0/float64(years)) - 1
	return rate
}

func Atof(input string) float64 {
	tmp, err := strconv.ParseFloat(input, 64)
	if err != nil {
		return 0
	} else {
		return tmp
	}
}

func Sum(input []float64) (result float64) {
	for _, row := range input {
		result += row
	}

	return result
}

func Mean(input []float64) float64 {
	return Sum(input) / float64(len(input))
}

func Variance(input []float64) (result float64) {
	var (
		meanValue = Mean(input)
	)

	for _, row := range input {
		result += math.Pow(row-meanValue, 2)
	}

	return result / float64(len(input)-1)
}

// 标准差
func StandardDeviation(input []float64) (result float64) {
	return math.Sqrt(Variance(input))
}

// 格式化
func RoundFloat(val float64, precision uint) float64 {
	ratio := math.Pow(10, float64(precision))
	return math.Round(val*ratio) / ratio
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

// 调整正负
func ToNegativeValues(input []float64) (output []float64) {
	for _, row := range input {
		row2 := -row

		// 零不做处理，否则会展示成 -0.0
		if row == 0 {
			row2 = 0
		}

		output = append(output, row2)
	}

	return
}

// 复合增长率
func IsCAGRBiggerThan(input []float64, minCAGR float64) bool {
	if len(input) < 2 {
		return false
	}

	cagr := math.Pow(input[0]/input[len(input)-1], 1.0/float64(len(input)))
	return cagr > minCAGR
}

func IsAllValuesZero(input []float64, maxValue float64) bool {
	for _, row := range input {
		if row != 0 {
			return false
		}
	}

	return true
}

// 所有值都小于目标
func IsAllValuesSmallerThan(input []float64, maxValue float64) bool {
	for _, row := range input {
		if row > maxValue {
			return false
		}
	}

	return len(input) > 0
}

// 所有值都大于目标
func IsAllValuesBiggerThan(input []float64, minValue float64) bool {
	for _, row := range input {
		if row < minValue {
			return false
		}
	}

	return len(input) > 0
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
