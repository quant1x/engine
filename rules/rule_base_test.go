package rules

import (
	"fmt"
	"gitee.com/quant1x/engine/models"
	"testing"
)

func Test_baseFilter(t *testing.T) {
	code := "601868"
	code = "sh600622"
	code = "sh601188"
	code = "sz002682"
	stockShots := models.BatchSnapShot([]string{code})
	if len(stockShots) > 0 {
		snapshot := stockShots[0]
		passed, failKind, err := Filter(snapshot)
		fmt.Println(passed, failKind, err)
	}
}
