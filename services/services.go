package services

import (
	"gitee.com/quant1x/engine/cache"
	"gitee.com/quant1x/engine/cachel5"
	"gitee.com/quant1x/gotdx"
	"gitee.com/quant1x/gotdx/trading"
	"gitee.com/quant1x/gox/logger"
)

const (
	// 定制任务初始化cron定位9点
	cronInit            = "0 9 * * *"
	cronDefaultInterval = "@every 10s"
)

const (
	barIndexUpdateSnapshot = 1
	barIndexRealtimeKLine  = 2
)

func init() {
	// 定时重置缓存
	err := Register("clean", cronInit, jobGlobalReset)
	if err != nil {
		logger.Fatal(err)
	}
	// 刷新快照
	err = Register("update_snapshot", "@every 2s", jobUpdateSnapshot)
	if err != nil {
		logger.Fatal(err)
	}
	// 实时更新K线
	err = Register("realtime_kline", "", jobRealtimeKLine)
	if err != nil {
		logger.Fatal(err)
	}
	// 更新全部
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

// 状态是否交易中
func isTrading(status int) bool {
	if status == trading.ExchangeTrading || status == trading.ExchangeCallAuction {
		return true
	}
	return false
}
