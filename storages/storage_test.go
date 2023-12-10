package storages

import (
	"fmt"
	"gitee.com/quant1x/engine/trader"
	"gitee.com/quant1x/gotdx/trading"
	"testing"
)

func TestOrderFlag(t *testing.T) {
	date := trading.LastTradeDate()
	code := "sh600178"
	direction := trader.BUY
	filename := getOrderFlagFilename(date, code, direction)
	fmt.Println(filename)
	err := Touch(filename)
	fmt.Println(err)
	ok := checkOrderStatus(date, code, direction)
	fmt.Println(ok)
}
