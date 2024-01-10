package factors

import (
	"fmt"
	"gitee.com/quant1x/engine/datasource/base"
	"testing"
)

func TestTransaction(t *testing.T) {
	code := "sh880941"
	date := "2024-01-09"
	list := base.Transaction(code, date)
	v := CountInflow(list, code, date)
	fmt.Printf("%+v\n", v)
}
