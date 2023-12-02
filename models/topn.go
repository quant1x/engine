package models

import "gitee.com/quant1x/engine/config"

var (
	CountDays int // 统计多少天, 由控制台传入数值
	CountTopN int // 统计前N, 由控制台传入数值
)

// AllStockTopN 最大输出多少只个股
func AllStockTopN() int {
	return config.EngineConfig.Order.TopN
}
