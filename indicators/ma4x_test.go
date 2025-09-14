package indicators

import (
	"fmt"
	"testing"

	"github.com/quant1x/engine/datasource/base"
	"github.com/quant1x/exchange"
	"github.com/quant1x/pandas"
)

func TestMA4X(t *testing.T) {
	code := "300781"
	date := "2024-06-25"
	code = exchange.CorrectSecurityCode(code)
	date = exchange.FixTradeDate(date)
	rows := base.CheckoutKLines(code, date)
	df := pandas.LoadStructs(rows)
	fmt.Println(df)
	df1 := MA4X(df, 5)
	fmt.Println(df1)
}
