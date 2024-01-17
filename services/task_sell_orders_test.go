package services

import (
	"fmt"
	"gitee.com/quant1x/engine/trader"
	"gitee.com/quant1x/exchange"
	"testing"
)

func TestTaskSell_getEarlierDate(t *testing.T) {
	v := getEarlierDate(1)
	fmt.Println(v)
}

func Test_checkoutCanSellStockList(t *testing.T) {
	positions, err := trader.QueryHolding()
	if err != nil {
		return
	}
	var holdings []string
	for _, position := range positions {
		if position.CanUseVolume < 1 {
			continue
		}
		stockCode := position.StockCode
		securityCode := exchange.CorrectSecurityCode(stockCode)
		holdings = append(holdings, securityCode)
	}
	v := CheckoutCanSellStockList(117, holdings)
	fmt.Println(v)
}

func Test_getHoldingDates(t *testing.T) {
	v := getHoldingDates(1)
	fmt.Println(v)
}
