package dfcf

import (
	"fmt"
	"gitee.com/quant1x/gotdx/trading"
	"testing"
)

func TestShareHolder(t *testing.T) {
	code := "sh600115"
	v := ShareHolder(code, trading.Today(), 2)
	fmt.Println(v)
}

func TestGetCacheShareHolder(t *testing.T) {
	code := "sh600105"
	v := GetCacheShareHolder(code, trading.Today(), 4)
	fmt.Println(v)
}
