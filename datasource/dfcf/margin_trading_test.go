package dfcf

import (
	"fmt"
	"testing"
)

func TestMarginTrading(t *testing.T) {
	date := "20240408"
	v, n, err := rawMarginTradingList(date, 2)
	fmt.Println(v, n, err)
}
