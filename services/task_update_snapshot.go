package services

import (
	"gitee.com/quant1x/engine/models"
	"gitee.com/quant1x/exchange"
	"gitee.com/quant1x/gox/logger"
	"gitee.com/quant1x/gox/runtime"
	"time"
)

// 任务 - 更新快照
func jobUpdateSnapshot() {
	tm := time.Now()
	updateInRealTime, status := exchange.CanUpdateInRealtime(tm)
	// 交易时间更新数据
	if updateInRealTime && (IsTrading(status) || exchange.CheckCallAuctionClose(tm)) {
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
	models.SyncAllSnapshots()
	logger.Infof("同步snapshot...OK")
}
