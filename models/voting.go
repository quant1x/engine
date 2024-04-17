package models

import (
	"gitee.com/quant1x/engine/config"
)

const (
	// CACHE_STRATEGY_PATH 策略文件存储路径
	CACHE_STRATEGY_PATH = "strategy"
)

var (
	CountDays int // 统计多少天, 由控制台传入数值
	CountTopN int // 统计前N, 由控制台传入数值
)

// AllStockTopN 最大输出多少只个股
func AllStockTopN() int {
	return config.TraderConfig().TopN
}
