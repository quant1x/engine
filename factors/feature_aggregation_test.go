package factors

import (
	"fmt"
	"testing"

	"gitee.com/quant1x/gotdx/securities"
)

func TestGetL5History(t *testing.T) {
	code := "sz000713"
	v := GetL5History(code)
	fmt.Println(v)
}

func TestMisc(t *testing.T) {
	code := "sh000001"
	v := GetL5Misc(code)
	fmt.Println(v)
}

func TestFilterL5Misc(t *testing.T) {
	rows := FilterL5Misc(func(v *Misc) bool {
		c1 := v.BullPower > v.BearPower
		//c2 := v.BullPower > 0 && v.BearPower != 0
		c2 := v.PowerTrendPeriod == 1
		return c1 && c2
	}, "20240205")
	for _, v := range rows {
		fmt.Println(v.Code, securities.GetStockName(v.Code))
	}
	fmt.Println("total:", len(rows))
}
