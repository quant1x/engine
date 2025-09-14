package services

import (
	"time"

	"github.com/quant1x/engine/cache"
	"github.com/quant1x/engine/factors"
	"github.com/quant1x/engine/storages"
	"github.com/quant1x/exchange"
)

var (
	// 非交易日每天更新一次
	lastUpdateTime = "22:00:00.000"
	// 交易日每天更新2次
	allDateUpdateTimes = []string{"15:10:00.000", lastUpdateTime}
)

// 任务 - 更新全部数据
func jobUpdateAll() {
	now := time.Now()
	tm := now.Format(exchange.CN_SERVERTIME_FORMAT)
	today := exchange.Today()
	lastDate := exchange.LastTradeDate()
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
		factors.SwitchDate(cache.DefaultCanReadDate())
		updateAll()
		doneUpdate(today, phase)
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
	storages.DataSetUpdate(barIndex, featureDate, plugins, cache.OpUpdate)
}

func updateAllFeatures(barIndex int, cacheDate, featureDate string) {
	// 1. 获取全部注册的数据集插件
	mask := cache.PluginMaskFeature
	plugins := cache.Plugins(mask)
	storages.FeaturesUpdate(&barIndex, cacheDate, featureDate, plugins, cache.OpUpdate)
}
