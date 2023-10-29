package storages

import (
	"fmt"
	"gitee.com/quant1x/gox/api"
	"testing"
)

func TestStockPool(t *testing.T) {
	filename := "./t1.csv"
	sp := StockPool{
		Status: 2,
	}
	sp.Status.Set(StrategyPassed, true)
	list := []StockPool{sp}
	api.SlicesToCsv(filename, list)

	list1 := []StockPool{}

	api.CsvToSlices(filename, &list1)
	sp = list1[0]
	fmt.Println(sp.Status.IsCancel())
}
