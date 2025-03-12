package dfcf

import (
	"fmt"
	"gitee.com/quant1x/pandas"
	"testing"
)

func TestIndividualStocksFundFlow(t *testing.T) {
	code := "sz000701"
	list := IndividualStocksFundFlow(code, "2025-03-12")
	df := pandas.LoadStructs(list)
	fmt.Println(df)
}
