package storages

import (
	"fmt"
	"gitee.com/quant1x/exchange"
	"testing"
)

func Test_checkOrderForBuy(t *testing.T) {
	list := getStockPoolFromCache()
	model := TestModel{}
	date := exchange.LastTradeDate()
	v := checkOrderForBuy(list, model, date)
	fmt.Println(v)
	saveStockPoolToCache(list)
}

func Test_strategyOrderIsFinished(t *testing.T) {
	v := strategyOrderIsFinished(TestModel{})
	fmt.Println(v)
}
