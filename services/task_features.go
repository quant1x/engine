package services

import (
	"gitee.com/quant1x/engine/cache"
	"gitee.com/quant1x/engine/cachel5"
	"gitee.com/quant1x/engine/storages"
	"gitee.com/quant1x/gotdx/trading"
	"gitee.com/quant1x/gox/logger"
	"time"
)

var (
	// 非交易日每天更新一次
	lastUpdateTime = "22:00:00.000"
	// 交易日每天更新2次
	allDateUpdateTimes = []string{"15:10:00.000", "18:10:00.000", lastUpdateTime}
)

// 任务 - 更新全部数据
func jobUpdateAll() {
	now := time.Now()
	tm := now.Format(trading.CN_SERVERTIME_FORMAT)
	today := trading.Today()
	lastDate := trading.LastTradeDate()
	bUpdated := false
	phase := ""
	if today == lastDate {
		for _, v := range allDateUpdateTimes {
			if tm >= v {
				phase = v
				bUpdated = checkUpdateState(today, phase)
				if bUpdated {
					break
				}
			}
		}
	} else {
		if tm >= lastUpdateTime {
			phase = lastUpdateTime
			bUpdated = checkUpdateState(today, phase)
		}
	}
	if bUpdated && len(phase) > 0 {
		cachel5.SwitchDate(cache.DefaultCanReadDate())
		updateAll()
		doneUpdate(today, phase)
	} else {
		logger.Infof("非全数据更新时段")
	}
}

func updateAll() {
	barIndex := 1
	currentDate := cache.DefaultCanUpdateDate()
	cacheDate, featureDate := cache.CorrectDate(currentDate)
	updateAllBaseData(barIndex, featureDate)
	updateAllFeatures(barIndex+1, cacheDate, featureDate)
}

func updateAllBaseData(barIndex int, featureDate string) {
	// 1. 获取全部注册的数据集插件
	mask := cache.PluginMaskBaseData
	plugins := cache.Plugins(mask)
	// 2. 执行操作
	storages.BaseDataUpdate(barIndex, featureDate, plugins, cache.OpUpdate)
}

func updateAllFeatures(barIndex int, cacheDate, featureDate string) {
	// 1. 获取全部注册的数据集插件
	mask := cache.PluginMaskFeature
	//dataSetList := flash.DataSetList()
	plugins := cache.Plugins(mask)
	storages.FeaturesUpdate(&barIndex, cacheDate, featureDate, plugins, cache.OpUpdate)
}
