package storages

import (
	"gitee.com/quant1x/engine/cache"
	"gitee.com/quant1x/engine/config"
	"gitee.com/quant1x/engine/models"
	"path/filepath"
)

var (
	traderConfig = config.TraderConfig()
)

const (
	StrategiesPath = "quant" // 策略结果数据文件存储路径
)

// GetResultCachePath 获取结果缓存路径
func GetResultCachePath() string {
	path := filepath.Join(cache.GetRootPath(), StrategiesPath)
	return path
}

// OutputStatistics 输出策略结果
func OutputStatistics(model models.Strategy, date string, v []models.Statistics) {
	tradeRule := config.GetStrategyParameterByCode(model.Code())
	if tradeRule == nil || !tradeRule.Enable() || tradeRule.Total == 0 {
		// 配置不存在, 或者规则无效
		return
	}
	topN := tradeRule.Total
	stockPoolMerge(model, date, v, topN)
}
