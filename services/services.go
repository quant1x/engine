package services

import (
	"gitee.com/quant1x/engine/cache"
	"gitee.com/quant1x/engine/cachel5"
	"gitee.com/quant1x/gotdx"
	"gitee.com/quant1x/gotdx/trading"
	"gitee.com/quant1x/gox/logger"
)

const (
	// CronTriggerInit 定制任务初始化cron定位9点
	CronTriggerInit = "0 9 * * *"
	// CronDefaultInterval 默认的执行频次
	CronDefaultInterval = "@every 10s"
	// CronTickInterval tick级别数据的更新频次
	CronTickInterval = "@every 1s"
)

const (
	barIndexUpdateSnapshot = 1
	barIndexRealtimeKLine  = 2
)

// 定时任务关键字
const (
	keyCronReset            = "clean"           // 定时清理重置数据状态
	keyCronRealTimeKLine    = "realtime_kline"  //  实时更新K线
	keyCronUpdateSnapshot   = "update_snapshot" // 更新快照
	keyCronUpdateAll        = "update_all"      // 更新全部数据, 包括基础数据和特征数据
	keyCronCookieCutterSell = "sell_117"        // 一刀切卖出
)

func init() {
	// 定时重置缓存
	err := Register(keyCronReset, CronTriggerInit, jobGlobalReset)
	if err != nil {
		logger.Fatal(err)
	}
	// 刷新快照
	err = Register(keyCronUpdateSnapshot, CronTickInterval, jobUpdateSnapshot)
	if err != nil {
		logger.Fatal(err)
	}
	// 实时更新K线
	err = Register(keyCronRealTimeKLine, CronDefaultInterval, jobRealtimeKLine)
	if err != nil {
		logger.Fatal(err)
	}

	// 更新全部
	err = Register(keyCronUpdateAll, CronDefaultInterval, jobUpdateAll)
	if err != nil {
		logger.Fatal(err)
	}

	// 一刀切卖出
	err = Register(keyCronCookieCutterSell, CronDefaultInterval, jobOneSizeFitsAllSales)
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

// IsTrading 状态是否交易中
func IsTrading(status int) bool {
	if status == trading.ExchangeTrading || status == trading.ExchangeCallAuction {
		return true
	}
	return false
}
