package commons_test

import (
	"fmt"
	"testing"

	"github.com/CaledoniaProject/gopkg-commons/commons"
)

func TestSize(t *testing.T) {
	fmt.Println(commons.FormatBytes(1000000000))
}
