package services

import (
	"github.com/quant1x/exchange"
	"github.com/quant1x/x/logger"
)

const (
	// CronTriggerNetwork 网络重置每天8点55分
	CronTriggerNetwork = "55 8 * * *"
	// CronTriggerInit 定制任务初始化cron定位9点
	CronTriggerInit = "0 9 * * *"
	// CronDefaultInterval 默认的执行频次
	CronDefaultInterval = "@every 10s"
	// CronTickInterval tick级别数据的更新频次
	CronTickInterval = "@every 1s"
	// 同步订单触发时间, 每交易日15点~23点的02分
	cronSyncOrdersInterval = "2 15-23 * * *"
	// cronMarginTrading 更新融资融券
	cronMarginTrading = "5 9 * * *"
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
	keyCronUpdateMisc       = "update_misc"     // 更新misc
	keyCronUpdateAll        = "update_all"      // 更新全部数据, 包括基础数据和特征数据
	keyCronCookieCutterSell = "sell_117"        // 一刀切卖出, one-size-fits-all
	keyCronSyncQmtOrder     = "sync_orders"     // 同步订单
	keyCronResetNetwork     = "reset_network"   // 重置网络
	keyCronMarginTrading    = "update_rzrq"     // 更新融资融券
)

func init() {
	err := Register(keyCronResetNetwork, CronTriggerNetwork, jobResetNetwork)
	if err != nil {
		logger.Fatal(err)
	}
	// 定时重置缓存
	err = Register(keyCronReset, CronTriggerInit, jobGlobalReset)
	if err != nil {
		logger.Fatal(err)
	}
	// 实时更新快照
	err = Register(keyCronUpdateSnapshot, CronTickInterval, jobUpdateSnapshot)
	if err != nil {
		logger.Fatal(err)
	}
	// 更新misc特征数据
	err = Register(keyCronUpdateMisc, CronTickInterval, jobUpdateMiscAndSnapshot)
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
	// 更新融资融券
	err = Register(keyCronMarginTrading, cronMarginTrading, jobUpdateMarginTrading)
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
