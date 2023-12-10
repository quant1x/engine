package trader

import (
	"fmt"
	"gitee.com/quant1x/engine/models"
	"gitee.com/quant1x/engine/strategies"
	"gitee.com/quant1x/gox/util/treemap"
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

	orderId, err := PlaceOrder(direction, model, securityCode, price, volume)
	fmt.Println(orderId, err)
}

func TestCalculateFundForStrategy(t *testing.T) {
	var model models.Strategy
	model = new(TestModel)
	fund := CalculateFundForStrategy(model)
	fmt.Println(fund)
}

type TestModel struct{}

func (TestModel) Code() models.ModelKind {
	return 81
}

func (TestModel) Name() string {
	//TODO implement me
	panic("implement me")
}

func (TestModel) OrderFlag() string {
	//TODO implement me
	panic("implement me")
}

func (TestModel) Filter(snapshot models.QuoteSnapshot) bool {
	//TODO implement me
	panic("implement me")
}

func (TestModel) Sort(snapshots []models.QuoteSnapshot) models.SortedStatus {
	//TODO implement me
	panic("implement me")
}

func (TestModel) Evaluate(securityCode string, result *treemap.Map) {
	//TODO implement me
	panic("implement me")
}
