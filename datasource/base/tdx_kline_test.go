package base

import (
	"fmt"
	"gitee.com/quant1x/exchange"
	"gitee.com/quant1x/pandas"
	"testing"
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
	code = exchange.CorrectSecurityCode(code)
	klines := UpdateAllBasicKLine(code)
	df := pandas.LoadStructs(klines)
	fmt.Println(df)
	_ = df
}
