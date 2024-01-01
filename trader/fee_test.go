package trader

import (
	"fmt"
	"testing"
)

func TestFundAllocate(t *testing.T) {
	traderParameter.ResetPositionRatio()
	fmt.Println(traderParameter)
}

func TestEvaluateFeeForBuy(t *testing.T) {
	code := "sh600178"
	price := 8.17

	v := EvaluateFeeForBuy(code, traderParameter.BuyAmountMax, price)
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
