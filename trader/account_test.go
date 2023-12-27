package trader

import (
	"fmt"
	"gitee.com/quant1x/engine/config"
	"testing"
)

func TestCalculateAvailableFund(t *testing.T) {
	id := 2
	tradeRule := config.GetStrategyParameterByCode(id)
	if tradeRule == nil {
		return
	}
	fmt.Println(tradeRule)
	fund := CalculateAvailableFund(tradeRule)
	fmt.Println(fund)
}
