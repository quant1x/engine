package trader

import (
	"fmt"
	"gitee.com/quant1x/engine/config"
	"testing"
)

func TestCalculateAvailableFund(t *testing.T) {
	id := 82
	tradeRule := config.GetTradeRule(id)
	if tradeRule == nil {
		return
	}
	fmt.Println(tradeRule)
	fund := CalculateAvailableFund(tradeRule)
	fmt.Println(fund)
}
