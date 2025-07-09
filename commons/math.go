package commons

import (
	"math"
	"strconv"
)

func Atof(input string) float64 {
	tmp, err := strconv.ParseFloat(input, 64)
	if err != nil {
		return 0
	} else {
		return tmp
	}
}

// 求和
func Sum(input []float64) (result float64) {
	for _, row := range input {
		result += row
	}

	return result
}

// 均值
func Mean(input []float64) float64 {
	return Sum(input) / float64(len(input))
}

// 方差
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

// 所有值都是零
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
