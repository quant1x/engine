package storages

import (
	"fmt"
	"gitee.com/quant1x/engine/cache"
	"gitee.com/quant1x/engine/config"
	"gitee.com/quant1x/engine/models"
	"gitee.com/quant1x/exchange"
	"gitee.com/quant1x/gox/api"
	"gitee.com/quant1x/pandas"
	"os"
)

var (
	traderConfig = config.TraderConfig()
)

const (
	StrategiesPath = "zero-sum" // 策略结果数据文件存储路径

)

// GetResultCachePath 获取结果缓存路径
func GetResultCachePath() string {
	path := fmt.Sprintf("%s/%s", cache.GetRootPath(), StrategiesPath)
	return path
}

// OutputStatistics 输出策略结果
func OutputStatistics(model models.Strategy, date string, v []models.Statistics) {
	df := pandas.LoadStructs(v)
	if df.Nrow() == 0 {
		return
	}
	tradeRule := config.GetStrategyParameterByCode(model.Code())
	if tradeRule == nil || !tradeRule.Enable() {
		// 配置不存在, 或者规则无效
		return
	}
	topN := tradeRule.Total
	orderFlag := model.OrderFlag()
	date = exchange.FixTradeDate(date, cache.FilenameDate)
	filename := fmt.Sprintf("%s/%s-%d.csv", GetResultCachePath(), date, topN)
	_ = df.WriteCSV(filename)
	updateTime, _ := api.ParseTime(v[0].UpdateTime)
	if exchange.CanUpdate(updateTime) {
		fnOrder := fmt.Sprintf("%s/%s-%s.csv", cache.GetQmtCachePath(), date, orderFlag)
		if !api.FileExist(fnOrder) {
			err := df.WriteCSV(fnOrder)
			if err != nil {
				fmt.Println(err)
				return
			}
			fnReady := fmt.Sprintf("%s/%s-%s.ready", cache.GetQmtCachePath(), date, orderFlag)
			file, err := os.Create(fnReady)
			if err != nil {
				fmt.Println(err)
				return
			}
			api.CloseQuietly(file)
		}
	}
	stockPoolMerge(model, date, v, topN)
}
