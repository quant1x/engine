package dfcf

import (
	"fmt"
	"gitee.com/quant1x/pandas"
	"testing"
)

func TestIndividualStocksFundFlow(t *testing.T) {
	code := "sh603260"
	list := IndividualStocksFundFlow(code, "2023-06-09")
	df := pandas.LoadStructs(list)
	fmt.Println(df)
}
