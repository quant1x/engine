package rules

import (
	"gitee.com/quant1x/engine/config"
	"gitee.com/quant1x/gox/api"
)

const (
	TenThousand = 1e4               // 万
	Million     = 100 * TenThousand // 百万
	Billion     = 100 * Million     // 1亿
)

var (
	RuleParameters config.RuleParameter
)

func init() {
	// 初始化配置
	rules := config.RuleConfig()
	// 加载规则参数
	_ = api.Copy(&RuleParameters, &rules)
	// 流通盘
	//RuleParameters.CapitalMin *= Billion
	//RuleParameters.CapitalMax *= Billion
	// 市值
	//RuleParameters.MarketCapMin *= Billion
	//RuleParameters.MarketCapMax *= Billion
	// 最大流出
	//RuleParameters.MaxReduceAmount *= TenThousand
}
