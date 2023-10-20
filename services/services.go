package services

import (
	"gitee.com/quant1x/engine/cache"
	"gitee.com/quant1x/engine/cachel5"
	"gitee.com/quant1x/engine/storages"
	"gitee.com/quant1x/gotdx"
	"gitee.com/quant1x/gotdx/trading"
	"gitee.com/quant1x/gox/cron"
	"gitee.com/quant1x/gox/logger"
	"gitee.com/quant1x/gox/signal"
	"time"
)

const (
	// 定制任务初始化cron定位9点
	cronInit = "0 9 * * *"
)

var (
	// 非交易日每天更新一次
	lastUpdateTime = "22:00:00.000"
	// 交易日每天更新2次
	allDateUpdateTimes = []string{"15:10:00.000", "18:10:00.000", lastUpdateTime}
)

func init() {
	err := Register("clean", cronInit, cleanStatFiles)
	if err != nil {
		logger.Fatal(err)
	}
	//err = Register("update_all", "", callbackUpdateAll)
	//if err != nil {
	//	logger.Fatal(err)
	//}
}

func cleanStatFiles() {
	logger.Info("清理过期的更新状态文件...")
	_ = cleanExpiredStateFiles()
	gotdx.ReOpen()
	cachel5.SwitchDate(cache.DefaultCanReadDate())
	logger.Info("清理过期的更新状态文件...OK")
}

// v1DaemonService 守护进程服务入口
func v1DaemonService() {
	job := cron.New()
	// 1. 数据初始化
	// 1.1. 定时清理过期的状态文件
	job.Start()
	_, err := job.AddFunc(cronInit, func() {
		logger.Info("清理过期的更新状态文件...")
		_ = cleanExpiredStateFiles()
		gotdx.ReOpen()
		cachel5.SwitchDate(cache.DefaultCanReadDate())
		logger.Info("清理过期的更新状态文件...OK")
	})
	if err != nil {
		logger.Fatal(err)
		return
	}
	// 2. 全部更新
	_, err = job.AddFuncWithSkipIfStillRunning("@every 10s", func() {
		callbackUpdateAll()
	})
	if err != nil {
		logger.Fatal(err)
		return
	}
	// 3. 盘中增量更新
	// 4. 阻塞
	interrupt := signal.Notify()
	select {
	case sig := <-interrupt:
		logger.Infof("interrupt: %s", sig.String())
		break
	}
}

// 更新全部数据
func callbackUpdateAll() {
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
