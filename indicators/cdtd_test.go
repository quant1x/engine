package indicators

import (
	"fmt"
	"testing"

	"github.com/quant1x/engine/datasource/base"
	"github.com/quant1x/exchange"
	"github.com/quant1x/pandas"
)

func TestCDTD(t *testing.T) {
	code := "002528.sz"
	code = "600839"
	date := "2025-02-14"
	code = exchange.CorrectSecurityCode(code)
	date = exchange.FixTradeDate(date)
	rows := base.CheckoutKLines(code, date)
	df := pandas.LoadStructs(rows)
	fmt.Println(df)
	df1 := CDTD(df)
	fmt.Println(df1)
}
