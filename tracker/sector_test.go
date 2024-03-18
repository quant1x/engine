package tracker

import (
	"fmt"
	"gitee.com/quant1x/gotdx/securities"
	"gitee.com/quant1x/pandas"
	"testing"
)

func Test_scanBlock(t *testing.T) {
	pbarIndex := 0
	data := scanSectorSnapshots(&pbarIndex, securities.BK_HANGYE, false)
	df := pandas.LoadStructs(data)
	fmt.Println(df)
}
