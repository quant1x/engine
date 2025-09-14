package services

import (
	"time"

	"github.com/quant1x/engine/models"
	"github.com/quant1x/exchange"
	"github.com/quant1x/x/logger"
	"github.com/quant1x/x/runtime"
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
	models.SyncAllSnapshots(nil)
	logger.Infof("同步snapshot...OK")
}
