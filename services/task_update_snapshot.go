package services

import (
	"gitee.com/quant1x/engine/models"
	"gitee.com/quant1x/exchange"
	"gitee.com/quant1x/gox/logger"
	"gitee.com/quant1x/gox/runtime"
)

// 任务 - 更新快照
func jobUpdateSnapshot() {
	updateInRealTime, status := exchange.CanUpdateInRealtime()
	// 交易时间更新数据
	if updateInRealTime && IsTrading(status) {
		realtimeUpdateSnapshot()
	} else {
		if runtime.Debug() {
			realtimeUpdateSnapshot()
		}
	}
}

// 更新快照
func realtimeUpdateSnapshot() {
	logger.Infof("同步snapshot...")
	models.SyncAllSnapshots(nil)
	logger.Infof("同步snapshot...OK")
}
