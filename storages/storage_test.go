package storages

import (
	"fmt"
	"gitee.com/quant1x/engine/factors"
	"gitee.com/quant1x/engine/models"
	"gitee.com/quant1x/engine/trader"
	"gitee.com/quant1x/gotdx/trading"
	"gitee.com/quant1x/gox/concurrent"
	"testing"
)

func TestOrderFlag(t *testing.T) {
	model := TestModel82{}
	date := trading.LastTradeDate()
	code := "sh600178"
	direction := trader.BUY
	filename := order_state_filename(date, model, code, direction)
	fmt.Println(filename)
	err := Touch(filename)
	fmt.Println(err)
	ok := CheckOrderState(date, model, code, direction)
	fmt.Println(ok)
}

func TestCountOrderFlag(t *testing.T) {
	model := TestModel82{}
	date := trading.LastTradeDate()
	direction := trader.BUY
	v := CountStrategyOrders(date, model, direction)
	fmt.Println(v)
}

func TestGetOrderDateFirstBuy(t *testing.T) {
	date := "2023-12-13"
	strategyName := "S82"
	direction := trader.BUY
	v := FetchListForFirstPurchase(date, strategyName, direction)
	fmt.Println(v)
}

type TestModel82 struct{}

func (TestModel82) Code() models.ModelKind {
	return 82
}

func (TestModel82) Name() string {
	//TODO implement me
	panic("implement me")
}

func (TestModel82) OrderFlag() string {
	//TODO implement me
	panic("implement me")
}

func (TestModel82) Filter(snapshot factors.QuoteSnapshot) error {
	//TODO implement me
	panic("implement me")
}

func (TestModel82) Sort(snapshots []factors.QuoteSnapshot) models.SortedStatus {
	//TODO implement me
	panic("implement me")
}

func (TestModel82) Evaluate(securityCode string, result *concurrent.TreeMap[string, models.ResultInfo]) {
	//TODO implement me
	panic("implement me")
}
