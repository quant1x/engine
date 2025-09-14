package storages

import (
	"fmt"
	"testing"

	"github.com/quant1x/engine/trader"
	"github.com/quant1x/exchange"
	"github.com/quant1x/x/api"
)

func TestOrderFlag(t *testing.T) {
	model := TestModel82{}
	date := exchange.LastTradeDate()
	code := "sh600178"
	direction := trader.BUY
	filename := order_state_filename(date, model, direction, code)
	fmt.Println(filename)
	err := api.Touch(filename)
	fmt.Println(err)
	ok := CheckOrderState(date, model, code, direction)
	fmt.Println(ok)
}

func TestCountOrderFlag(t *testing.T) {
	model := TestModel82{}
	date := exchange.LastTradeDate()
	direction := trader.BUY
	v := CountStrategyOrders(date, model, direction)
	fmt.Println(v)
}

func TestGetOrderDateFirstBuy(t *testing.T) {
	date := "2024-05-17"
	strategyName := "S0"
	direction := trader.BUY
	v := FetchListForFirstPurchase(date, strategyName, direction)
	fmt.Println(v)
}
