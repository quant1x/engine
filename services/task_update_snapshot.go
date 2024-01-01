package services

import (
	"gitee.com/quant1x/engine/models"
	"gitee.com/quant1x/gotdx/trading"
	"gitee.com/quant1x/gox/logger"
	"gitee.com/quant1x/gox/runtime"
	"runtime/debug"
)

// 任务 - 更新快照
func jobUpdateSnapshot() {
	//funcName, _, _ := runtime.Caller()
	updateInRealTime, status := trading.CanUpdateInRealtime()
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
	defer func() {
		if err := recover(); err != nil {
			s := string(debug.Stack())
			logger.Errorf("err=%v, stack=%s", err, s)
		}
	}()
	models.SyncAllSnapshots(nil)
}
