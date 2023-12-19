package base

import (
	"fmt"
	"testing"
)

func TestFundFlow(t *testing.T) {
	code := "sh600977"
	code = "sh600058"
	code = "sz002175"
	code = "sz002043"
	df := FundFlow(code)
	fmt.Println(df)
}
