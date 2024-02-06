package base

import (
	"fmt"
	"testing"
)

func TestGetBeginDateOfHistoricalTradingData(t *testing.T) {
	v := GetBeginDateOfHistoricalTradingData()
	fmt.Println(v)
}
