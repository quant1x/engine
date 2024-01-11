package services

import (
	"gitee.com/quant1x/exchange"
	"gitee.com/quant1x/gox/logger"
)

const (
	// CronTriggerInit 定制任务初始化cron定位9点
	CronTriggerInit = "0 9 * * *"
	// CronDefaultInterval 默认的执行频次
	CronDefaultInterval = "@every 10s"
	// CronTickInterval tick级别数据的更新频次
	CronTickInterval = "@every 1s"
	// 同步订单触发时间, 每交易日15点02分
	cronSyncOrdersInterval = "2 15 * * *"
)

const (
	barIndexUpdateSnapshot            = 1
	barIndexRealtimeKLine             = 2
	barIndexUpdateExchangeAndSnapshot = 3
)

// 定时任务关键字
const (
	keyCronReset            = "global_reset"    // 全局重置
	keyCronRealTimeKLine    = "realtime_kline"  // 实时更新K线
	keyCronUpdateSnapshot   = "update_snapshot" // 更新快照
	keyCronUpdateExchange   = "update_exchange" // 更新exchange
	keyCronUpdateAll        = "update_all"      // 更新全部数据, 包括基础数据和特征数据
	keyCronCookieCutterSell = "sell_117"        // 一刀切卖出, one-size-fits-all
	keyCronSyncQmtOrder     = "sync_orders"     // 同步订单
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
	// 更新快照
	err = Register(keyCronUpdateExchange, CronTickInterval, jobUpdateExchangeAndSnapshot)
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
	// 同步QMT订单
	err = Register(keyCronSyncQmtOrder, cronSyncOrdersInterval, jobSyncTraderOrders)
	if err != nil {
		logger.Fatal(err)
	}
}

// IsTrading 状态是否交易中
func IsTrading(status int) bool {
	if status == exchange.ExchangeTrading || status == exchange.ExchangeCallAuction {
		return true
	}
	return false
}
