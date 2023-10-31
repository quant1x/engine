package rules

import (
	"gitee.com/quant1x/engine/cache"
	"gitee.com/quant1x/gox/api"
)

const (
	TenThousand = 1e4               // 万
	Million     = 100 * TenThousand // 百万
	Billion     = 100 * Million     // 1亿
)

var (
	RuleParameters  = cache.RuleParameter{}
	OrderParameters = cache.OrderParameter{}
)

func init() {
	// 初始化配置
	//cfg := EngineConfig.ReadConfig()
	// 加载规则参数
	_ = api.Copy(&RuleParameters, &cache.EngineConfig.Rules)
	RuleParameters.CapitalMin *= Billion
	RuleParameters.CapitalMax *= Billion
	RuleParameters.MaxReduceAmount *= TenThousand

	// 加载订单参数
	_ = api.Copy(&OrderParameters, &cache.EngineConfig.Order)
}
