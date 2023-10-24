package strategies

import "gitee.com/quant1x/engine/models"

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
	return true
}
