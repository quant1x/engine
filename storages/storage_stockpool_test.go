package storages

import (
	"fmt"
	"gitee.com/quant1x/engine/models"
	"gitee.com/quant1x/gotdx/trading"
	"gitee.com/quant1x/gox/util/treemap"
	"testing"
)

type TestModel struct{}

func (m TestModel) Code() models.ModelKind {
	return 81
}

func (m TestModel) Name() string {
	return "81号策略"
}

func (m TestModel) OrderFlag() string {
	return models.OrderFlagTick
}

func (m TestModel) Filter(snapshot models.QuoteSnapshot) bool {
	//TODO implement me
	panic("implement me")
}

func (m TestModel) Sort(snapshots []models.QuoteSnapshot) models.SortedStatus {
	//TODO implement me
	panic("implement me")
}

func (m TestModel) Evaluate(securityCode string, result *treemap.Map) {
	//TODO implement me
	panic("implement me")
}

func Test_strategyOrderIsFinished(t *testing.T) {
	v := strategyOrderIsFinished(TestModel{})
	fmt.Println(v)
}

func Test_checkOrderForBuy(t *testing.T) {
	list := GetStockPool()
	model := TestModel{}
	date := trading.LastTradeDate()
	v := checkOrderForBuy(list, model, date)
	fmt.Println(v)
}
