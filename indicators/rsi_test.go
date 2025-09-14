package indicators

import (
	"fmt"
	"testing"

	"github.com/quant1x/engine/datasource/base"
	"github.com/quant1x/exchange"
	"github.com/quant1x/pandas"
	. "github.com/quant1x/pandas/formula"
)

func TestRSI(t *testing.T) {
	code := "300781"
	date := "2024-06-25"
	code = exchange.CorrectSecurityCode(code)
	date = exchange.FixTradeDate(date)
	rows := base.CheckoutKLines(code, date)
	df := pandas.LoadStructs(rows)
	fmt.Println(df)
	df1 := RSI(df, 6, 12, 24)
	fmt.Println(df1)
}

func TestRSI2(t *testing.T) {
	VALUES := []float64{44.34, 44.09, 44.15, 43.61, 44.33, 44.83, 45.10, 45.42, 45.84,
		46.08, 45.89, 46.03, 45.61, 46.28, 46.28, 46.00, 46.03, 46.41,
		46.22, 45.64, 46.21, 46.25, 45.71, 46.45, 45.78, 45.35, 44.03,
		44.18, 44.22, 44.57, 43.42, 42.66, 43.13}
	CLOSE := pandas.ToSeries(VALUES...)
	//	LC:=REF(CLOSE,1);
	LC := REF(CLOSE, 1)
	cls := CLOSE.Sub(LC)
	result := SMA(MAX(cls, 0), 14, 1).Div(SMA(ABS(cls), 14, 1)).Mul(100)
	fmt.Println(result)
}
