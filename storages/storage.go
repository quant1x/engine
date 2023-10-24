package storages

import (
	"fmt"
	"gitee.com/quant1x/engine/cache"
	"gitee.com/quant1x/engine/models"
	"gitee.com/quant1x/gotdx/trading"
	"gitee.com/quant1x/gox/api"
	"gitee.com/quant1x/pandas"
	"os"
)

const (
	StrategiesPath = "zero-sum" // 策略结果数据文件存储路径
	OrdersPath     = "qmt"      // QMT订单缓存路径
)

// GetResultCachePath 获取结果缓存路径
func GetResultCachePath() string {
	path := fmt.Sprintf("%s/%s", cache.GetRootPath(), StrategiesPath)
	return path
}

func getQmtCachePath() string {
	path := fmt.Sprintf("%s/%s", cache.GetRootPath(), OrdersPath)
	return path
}

func OutputStatistics(orderFlag string, top int, date string, v []models.Statistics) {
	df := pandas.LoadStructs(v)
	if df.Nrow() == 0 {
		return
	}
	date = trading.FixTradeDate(date, cache.FilenameDate)
	filename := fmt.Sprintf("%s/%s-%d.csv", GetResultCachePath(), date, top)
	_ = df.WriteCSV(filename)
	updateTime, _ := api.ParseTime(v[0].UpdateTime)
	if trading.CanUpdate(updateTime) {
		fnOrder := fmt.Sprintf("%s/%s-%s.csv", getQmtCachePath(), date, orderFlag)
		if !api.FileExist(fnOrder) {
			err := df.WriteCSV(fnOrder)
			if err != nil {
				fmt.Println(err)
				return
			}
			fnReady := fmt.Sprintf("%s/%s-%s.ready", getQmtCachePath(), date, orderFlag)
			file, err := os.Create(fnReady)
			if err != nil {
				fmt.Println(err)
				return
			}
			api.CloseQuietly(file)
		}
	}

}
