package tracker

import (
	"fmt"
	"gitee.com/quant1x/engine/config"
	"gitee.com/quant1x/engine/models"
)

// CheckStrategy 检查当前交易日中个股在策略中的执行情况
func CheckStrategy(strategyCode int, securityCode string) {
	snapshot := models.GetTick(securityCode)
	if snapshot == nil {
		fmt.Println("快照获取失败")
		return
	}
	strategyParameter := config.GetStrategyParameterByCode(strategyCode)
	if strategyParameter == nil {
		fmt.Printf("找不到%d号策略的配置\n", strategyCode)
		return
	}
	model, err := models.CheckoutStrategy(strategyCode)
	if err != nil {
		fmt.Println(err)
		return
	}
	v := model.Filter(strategyParameter.Rules, *snapshot)
	fmt.Println(v)
}
