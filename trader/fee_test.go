package trader

import (
	"fmt"
	"testing"
)

func TestFundAllocate(t *testing.T) {
	traderConfig.ResetPositionRatio()
	fmt.Println(traderConfig)
}

func TestEvaluateFeeForBuy(t *testing.T) {
	code := "sh600178"
	price := 8.17

	v := EvaluateFeeForBuy(code, traderConfig.BuyAmountMax, price)
	fmt.Println(v)
	v.log()
}

func TestEvaluateFeeForSell(t *testing.T) {
	code := "sh600178"
	price := 15.24
	volume := 5000

	v := EvaluateFeeForSell(code, price, volume)
	fmt.Println(v)
	v.log()
}
