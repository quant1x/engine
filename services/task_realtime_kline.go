package services

import (
	"gitee.com/quant1x/data/exchange"
	"gitee.com/quant1x/engine/datasource/base"
	"gitee.com/quant1x/engine/market"
	"gitee.com/quant1x/engine/models"
	"gitee.com/quant1x/gox/coroutine"
	"gitee.com/quant1x/gox/logger"
	"gitee.com/quant1x/gox/progressbar"
	"gitee.com/quant1x/gox/runtime"
)

// 任务 - 实时更新K线
func jobRealtimeKLine() {
	funcName := "jobRealtimeKLine"
	updateInRealTime, status := exchange.CanUpdateInRealtime()
	// 14:30:00~15:01:00之间更新数据
	if updateInRealTime && IsTrading(status) {
		realtimeUpdateOfKLine()
	} else {
		if runtime.Debug() {
			realtimeUpdateOfKLine()
		} else {
			logger.Infof("%s, 非尾盘交易时段: %d", funcName, status)
		}
	}
}

// 更新K线
func realtimeUpdateOfKLine() {
	defer runtime.IgnorePanic("")
	barIndex := barIndexRealtimeKLine
	allCodes := market.GetCodeList()
	wg := coroutine.NewRollingWaitGroup(5)
	bar := progressbar.NewBar(barIndex, "执行[实时更新K线]", len(allCodes))
	for _, code := range allCodes {
		updateKLine := func(waitGroup *coroutine.RollingWaitGroup, securityCode string) {
			defer waitGroup.Done()
			bar.Add(1)
			snapshot := models.GetTickFromMemory(securityCode)
			if snapshot != nil {
				base.BasicKLineForSnapshot(*snapshot)
			}
		}
		wg.Add(1)
		go updateKLine(wg, code)
	}
	wg.Wait()
}
