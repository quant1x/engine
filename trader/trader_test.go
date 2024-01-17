package trader

import (
	"fmt"
	"gitee.com/quant1x/engine/models"
	"gitee.com/quant1x/engine/strategies"
	"testing"
)

func TestQueryAccount(t *testing.T) {
	info, err := QueryAccount()
	fmt.Println(info, err)
}

func TestQueryHolding(t *testing.T) {
	info, err := QueryHolding()
	fmt.Println(info, err)
}

func TestQueryOrders(t *testing.T) {
	info, err := QueryOrders()
	fmt.Println(info, err)
}

func TestTradeCancelOrder(t *testing.T) {
	orderId := 1086140321
	err := CancelOrder(orderId)
	fmt.Println(err)
}

func TestTradePlaceOrder(t *testing.T) {
	direction := BUY
	model := strategies.ModelNo1{}
	securityCode := "sh600178"
	price := 13.68
	volume := 100

	orderId, err := PlaceOrder(direction, model, securityCode, FIX_PRICE, price, volume)
	fmt.Println(orderId, err)
}

func TestCalculateFundForStrategy(t *testing.T) {
	var model models.Strategy
	model = new(TestModel)
	fund := CalculateFundForStrategy(model)
	fmt.Println(fund)
}
