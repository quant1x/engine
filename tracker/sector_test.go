package tracker

import (
	"fmt"
	"testing"

	"gitee.com/quant1x/gotdx/securities"
	"gitee.com/quant1x/pandas"
)

func Test_scanBlock(t *testing.T) {
	pbarIndex := 0
	data := scanSectorSnapshots(&pbarIndex, securities.BK_HANGYE, false)
	df := pandas.LoadStructs(data)
	fmt.Println(df)
}
