package storages

import (
	"fmt"
	"gitee.com/quant1x/engine/models"
	"gitee.com/quant1x/engine/trader"
	"gitee.com/quant1x/gotdx/trading"
	"gitee.com/quant1x/gox/util/treemap"
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

func (TestModel82) Filter(snapshot models.QuoteSnapshot) bool {
	//TODO implement me
	panic("implement me")
}

func (TestModel82) Sort(snapshots []models.QuoteSnapshot) models.SortedStatus {
	//TODO implement me
	panic("implement me")
}

func (TestModel82) Evaluate(securityCode string, result *treemap.Map) {
	//TODO implement me
	panic("implement me")
}
