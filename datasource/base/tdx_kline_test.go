package base

import (
	"fmt"
	"testing"

	"gitee.com/quant1x/exchange"
	"gitee.com/quant1x/pandas"
)

func TestGetAllBasicKLine(t *testing.T) {
	code := "sh000001"
	code = "sz002043"
	code = "sz000567"
	code = "sz300580"
	code = "sz301129"
	code = "sz002669"
	code = "600256"
	code = "000063"
	code = "688981"
	code = "002857"
	code = "603230"
	code = "880866"
	code = "002350"
	code = exchange.CorrectSecurityCode(code)
	klines := UpdateAllBasicKLine(code)
	df := pandas.LoadStructs(klines)
	fmt.Println(df)
	_ = df
}

func Test_getYearDay(t *testing.T) {
	date := "2025-08-08"
	s, e := getYearDay(date)
	fmt.Println(s, e)
}

func TestConvertKlinesTrading(t *testing.T) {
	code := "sh000001"
	date := "2025-08-07"
	klines := CheckoutKLines(code, date)
	df := pandas.LoadStructs(klines)
	fmt.Println(df)
	tmpLines := ConvertKlinesTrading(klines, "q")
	df = pandas.LoadStructs(tmpLines)
	fmt.Println(df)
}
