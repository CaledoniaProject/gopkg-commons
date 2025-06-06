package commons_test

import (
	"testing"

	"github.com/CaledoniaProject/gopkg-commons/commons"
)

func TestDateDiff(t *testing.T) {
	days, err := commons.DaysToNow("2055-01-01 00:00:00")
	if err != nil {
		t.Fatalf("date failed: %v", err)
	} else {
		t.Logf("days: %v", days)
	}
}
