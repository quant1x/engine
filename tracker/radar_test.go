package tracker

import (
	"fmt"
	"gitee.com/quant1x/engine/config"
	"testing"
)

func TestConfig(t *testing.T) {
	strategyCode := 82
	rule := config.GetTradeRule(strategyCode)
	fmt.Println(rule)
	fmt.Println(rule.StockList())
}
