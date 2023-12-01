package strategies

import (
	"gitee.com/quant1x/engine/config"
	"gitee.com/quant1x/engine/models"
	"gitee.com/quant1x/engine/rules"
)

var (
	CountDays int // 统计多少天
	CountTopN int // 统计前N
)

// AllStockTopN 最大输出多少只个股
func AllStockTopN() int {
	return config.EngineConfig.Order.TopN
}

// RuleFilter 过滤条件
func RuleFilter(snapshot models.QuoteSnapshot) bool {
	passed, failed, err := rules.Filter(snapshot)
	if failed != rules.Pass {
		return false
	}
	_ = passed
	_ = err
	return true
}
