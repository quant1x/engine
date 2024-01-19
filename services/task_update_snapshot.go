package services

import (
	"gitee.com/quant1x/engine/models"
	"gitee.com/quant1x/exchange"
	"gitee.com/quant1x/gox/logger"
	"gitee.com/quant1x/gox/runtime"
)

// 任务 - 更新快照
func jobUpdateSnapshot() {
	//funcName, _, _ := runtime.Caller()
	updateInRealTime, status := exchange.CanUpdateInRealtime()
	// 14:30:00~15:01:00之间更新数据
	if updateInRealTime && IsTrading(status) {
		realtimeUpdateSnapshot()
	} else {
		if runtime.Debug() {
			realtimeUpdateSnapshot()
		}
		//else {
		//	logger.Infof("%s, 非交易时段: %d", funcName, status)
		//}
	}
}

// 更新快照
func realtimeUpdateSnapshot() {
	defer runtime.IgnorePanic()
	logger.Infof("同步snapshot...")
	models.SyncAllSnapshots(nil)
	logger.Infof("同步snapshot...OK")
}
