package market

import (
	"fmt"
	"testing"
)

func TestGetCodeList(t *testing.T) {
	codes := GetCodeList()
	fmt.Println(len(codes))
	codes = GetStockCodeList()
	fmt.Println(len(codes))
}

func TestPriceLimit(t *testing.T) {
	code := "sh603598"
	lastClose := 19.00

	up, down := PriceLimit(code, lastClose)
	fmt.Println(up, down)
}
