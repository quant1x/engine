package services

import (
	"gitee.com/quant1x/engine/datasets/base"
	"gitee.com/quant1x/engine/market"
	"gitee.com/quant1x/gotdx/quotes"
	"gitee.com/quant1x/gotdx/trading"
	"gitee.com/quant1x/gox/logger"
	"runtime/debug"
	"sync"
	"time"
)

// 任务 - 实时更新K线
func jobRealtimeKLine() {
	updateInRealTime, status := trading.CanUpdateInRealtime()
	// 14:30:00~15:01:00之间更新数据
	if updateInRealTime && isTrading(status) {
		realtimeUpdateOfKLine()
	} else {
		logger.Infof("非尾盘交易时段: %d", status)
	}
}

// 更新K线
func realtimeUpdateOfKLine() {
	mainStart := time.Now()
	defer func() {
		if err := recover(); err != nil {
			s := string(debug.Stack())
			logger.Errorf("err=%v, stack=%s", err, s)
		}
		elapsedTime := time.Since(mainStart) / time.Millisecond
		logger.Infof("总耗时: %.3fs", float64(elapsedTime)/1000)
	}()
	allCodes := market.GetCodeList()
	count := len(allCodes)
	var wg sync.WaitGroup
	for start := 0; start < count; start += quotes.TDX_SECURITY_QUOTES_MAX {
		length := count - start
		if length >= quotes.TDX_SECURITY_QUOTES_MAX {
			length = quotes.TDX_SECURITY_QUOTES_MAX
		}
		subCodes := []string{}
		for i := 0; i < length; i++ {
			securityCode := allCodes[start+i]
			subCodes = append(subCodes, securityCode)
		}
		if len(subCodes) == 0 {
			continue
		}
		updateKLine := func(waitGroup *sync.WaitGroup, codes []string) {
			waitGroup.Done()
			for i := 0; i < quotes.DefaultRetryTimes; i++ {
				err := base.BatchRealtimeBasicKLine(codes)
				if err != nil {
					logger.Errorf("ZS: 网络异常: %+v, 重试: %d", err, i+1)
					time.Sleep(time.Second * 1)
					continue
				}
				break
			}
		}
		wg.Add(1)
		go updateKLine(&wg, subCodes)
	}
	wg.Wait()
}
