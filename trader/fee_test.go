package trader

import (
	"fmt"
	"testing"
)

func TestEvaluateFeeForBuy(t *testing.T) {
	code := "sh600178"
	price := 8.17

	v := EvaluateFeeForBuy(code, traderConfig.Head.FeeMax, price)
	fmt.Println(v)

}
