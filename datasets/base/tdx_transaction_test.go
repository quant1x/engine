package base

import (
	"fmt"
	"testing"
)

func TestTransaction(t *testing.T) {
	code := "sz000421"
	code = "sh000001"
	date := "2023-08-16"
	list := Transaction(code, date)
	v := CountInflow(list, code, date)
	fmt.Println(v)
}
