package services

import (
	"gitee.com/quant1x/engine/datasets/base"
	"gitee.com/quant1x/engine/market"
	"gitee.com/quant1x/gotdx/quotes"
	"gitee.com/quant1x/gox/logger"
	"runtime/debug"
	"time"
)

// 实时更新K线
func jobRealtimeKLine() {
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
		for i := 0; i < quotes.DefaultRetryTimes; i++ {
			err := base.BatchRealtimeBasicKLine(subCodes)
			if err != nil {
				logger.Errorf("ZS: 网络异常: %+v, 重试: %d", err, i+1)
				continue
			}
			break
		}
	}
}
