package services

import (
	"gitee.com/quant1x/engine/cache"
	"gitee.com/quant1x/engine/cachel5"
	"gitee.com/quant1x/engine/datasets/base"
	"gitee.com/quant1x/engine/market"
	"gitee.com/quant1x/engine/storages"
	"gitee.com/quant1x/gotdx"
	"gitee.com/quant1x/gotdx/quotes"
	"gitee.com/quant1x/gotdx/trading"
	"gitee.com/quant1x/gox/logger"
	"runtime/debug"
	"sync"
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
	// 定时重置缓存
	err := Register("clean", cronInit, jobGlobalReset)
	if err != nil {
		logger.Fatal(err)
	}
	err = Register("realtime_kline", "", jobRealtimeKLine)
	if err != nil {
		logger.Fatal(err)
	}
	err = Register("update_all", "", jobUpdateAll)
	if err != nil {
		logger.Fatal(err)
	}
}

// 任务 - 交易日数据缓存重置
func jobGlobalReset() {
	logger.Info("清理过期的更新状态文件...")
	_ = cleanExpiredStateFiles()
	gotdx.ReOpen()
	cachel5.SwitchDate(cache.DefaultCanReadDate())
	logger.Info("清理过期的更新状态文件...OK")
}

// 任务 - 实时更新K线
func jobRealtimeKLine() {
	now := time.Now()
	updateInRealTime, status := trading.CanUpdateInRealtime()
	// 14:30:00~15:01:00之间更新数据
	if updateInRealTime && trading.CheckCallAuctionTail(now) {
		realtimeUpdateOfKLine()
	} else {
		logger.Infof("非尾盘交易时段: %d", status)
	}
}

// 更新K线
func realtimeUpdateOfKLine() {
	mainStart := time.Now()
	defer func() {
		if err := recover(); err != nil {
			s := string(debug.Stack())
			logger.Errorf("err=%v, stack=%s", err, s)
		}
		elapsedTime := time.Since(mainStart) / time.Millisecond
		logger.Infof("总耗时: %.3fs", float64(elapsedTime)/1000)
	}()
	allCodes := market.GetCodeList()
	count := len(allCodes)
	var wg sync.WaitGroup
	for start := 0; start < count; start += quotes.TDX_SECURITY_QUOTES_MAX {
		length := count - start
		if length >= quotes.TDX_SECURITY_QUOTES_MAX {
			length = quotes.TDX_SECURITY_QUOTES_MAX
		}
		subCodes := []string{}
		for i := 0; i < length; i++ {
			securityCode := allCodes[start+i]
			subCodes = append(subCodes, securityCode)
		}
		if len(subCodes) == 0 {
			continue
		}
		updateKLine := func(waitGroup *sync.WaitGroup, codes []string) {
			waitGroup.Done()
			for i := 0; i < quotes.DefaultRetryTimes; i++ {
				err := base.BatchRealtimeBasicKLine(codes)
				if err != nil {
					logger.Errorf("ZS: 网络异常: %+v, 重试: %d", err, i+1)
					time.Sleep(time.Second * 1)
					continue
				}
				break
			}
		}
		wg.Add(1)
		go updateKLine(&wg, subCodes)
	}
	wg.Wait()
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
