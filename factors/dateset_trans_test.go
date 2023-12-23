package factors

import (
	"fmt"
	"gitee.com/quant1x/engine/datasource/base"
	"testing"
)

func TestTransaction(t *testing.T) {
	code := "sz000421"
	code = "sh000001"
	date := "2023-12-22"
	list := base.Transaction(code, date)
	v := CountInflow(list, code, date)
	fmt.Println(v)
}
