package indicators

import (
	"fmt"
	"testing"

	"gitee.com/quant1x/engine/datasource/base"
	"gitee.com/quant1x/exchange"
	"gitee.com/quant1x/pandas"
)

func TestPlatform(t *testing.T) {
	code := "600703.sh"
	code = "603789.sh"
	code = "sz000506"
	code = "sh603367"
	//code = "sz002275"
	code = "sz002665"
	code = "sz002528"
	//code = "sz000892"
	//code = "sz000905"
	code = "sh600641"
	code = "sh688031"
	code = "sz000988"
	code = "sh600105"
	code = "sz002292"
	code = "sh600354"
	code = "sh605577"
	code = "sh688662"
	code = "sz300678"
	code = "sh605162"
	code = "sz002992"
	date := "2024-06-25"
	code = exchange.CorrectSecurityCode(code)
	date = exchange.FixTradeDate(date)
	rows := base.CheckoutKLines(code, date)
	df := pandas.LoadStructs(rows)
	fmt.Println(df)
	df1 := Platform(df)
	fmt.Println(df1)
	_ = df1.WriteCSV("t02.csv")
}
