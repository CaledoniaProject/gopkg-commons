package commons_test

import (
	"testing"

	"github.com/CaledoniaProject/gopkg-commons/commons"
)

func TestMath(t *testing.T) {
	t.Logf("dev: %f", commons.StandardDeviation([]float64{79, 84, 84, 88, 92, 93, 94, 97, 98, 99, 100, 101, 101, 102, 102, 108, 110, 113, 118, 125}))
}

func TestValuesBigger(t *testing.T) {
	t.Logf("result: %v", commons.IsAllValuesBiggerThan([]float64{100, 110, 200}, 10))
	t.Logf("result: %v", commons.IsAllValuesBiggerThan([]float64{1, 110, 200}, 10))

	t.Logf("result: %v", commons.IsAllYoYBiggerThan([]float64{200, 150, 100}, 10))
	t.Logf("result: %v", commons.IsAllYoYBiggerThan([]float64{200, 350, 100}, 10))
}
