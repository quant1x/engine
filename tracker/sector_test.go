package tracker

import (
	"fmt"
	"gitee.com/quant1x/gotdx/securities"
	"gitee.com/quant1x/pandas"
	"testing"
)

func TestCheckBlock(t *testing.T) {
	pbarIndex := 0
	data := v1TopBlock(&pbarIndex)
	df := pandas.LoadStructs(data)
	fmt.Println(df)
}

func Test_scanBlock(t *testing.T) {
	pbarIndex := 0
	data := scanSectorSnapshots(&pbarIndex, securities.BK_HANGYE)
	df := pandas.LoadStructs(data)
	fmt.Println(df)
}
