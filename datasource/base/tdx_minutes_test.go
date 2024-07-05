package base

import (
	"fmt"
	"gitee.com/quant1x/exchange"
	"gitee.com/quant1x/pandas"
	"testing"
)

func TestGetMinutes(t *testing.T) {
	code := "sh000001"
	code = "sh510050"
	code = "sh600105"
	code = "880948"
	code = "603228"
	date := "2023-09-28"
	date = "2024-07-02"
	code = exchange.CorrectSecurityCode(code)
	date = exchange.FixTradeDate(date)
	v := GetMinutes(code, date)
	df := pandas.LoadStructs(v)
	fmt.Println(df)
}
