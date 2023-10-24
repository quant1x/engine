package strategies

import (
	"gitee.com/quant1x/engine/models"
	"gitee.com/quant1x/engine/rules"
)

var (
	CountDays int // 统计多少天
	CountTopN int // 统计前N
)

// AllStockTopN 最大输出多少只个股
func AllStockTopN() int {
	return globalOrderRules.TopN
}

// RuleFilter 过滤条件
func RuleFilter(snapshot models.QuoteSnapshot) bool {
	passed, failed := rules.Each(snapshot)
	if failed != rules.Pass {
		return false
	}
	_ = passed
	return true
}
