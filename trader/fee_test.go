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

func TestEvaluatePriceForSell(t *testing.T) {
	fixedYield := 0.03
	code := "sh600178"
	price := 15.24
	volume := 5000
	baseAmount := price * float64(volume)
	fmt.Println(baseAmount)
	v := EvaluatePriceForSell(code, price, volume, 0.03)
	fmt.Println(v, v.MarketValue/baseAmount >= (1+fixedYield))
	v.log()
}
