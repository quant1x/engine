package market

import (
	"fmt"
	"strings"
	"testing"

	"github.com/quant1x/exchange"
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

func TestGetQmtCodeList(t *testing.T) {
	batchMax := 1000
	codes := GetStockCodeList()
	codes = codes[0:batchMax]
	var list []string
	for _, v := range codes {
		_, mflag, symbol := exchange.DetectMarket(v)
		securityCode := fmt.Sprintf("%s.%s", symbol, strings.ToUpper(mflag))
		list = append(list, securityCode)
	}
	fmt.Println(strings.Join(list, ","))
	_ = batchMax
}
