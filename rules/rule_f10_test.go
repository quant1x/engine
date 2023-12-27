package rules

import (
	"fmt"
	"gitee.com/quant1x/engine/config"
	"gitee.com/quant1x/engine/models"
	"testing"
)

func Test_baseFilter(t *testing.T) {
	code := "601868"
	code = "sh600622"
	code = "sh601188"
	code = "sz002682"
	snapshot := models.GetTick(code)
	strategyParametet := config.GetStrategyParameterByCode(0)
	passed, failKind, err := Filter(strategyParametet.Rules, *snapshot)
	fmt.Println(passed, failKind, err)

}
