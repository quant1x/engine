package factors

import (
	"gitee.com/quant1x/engine/cache"
	"gitee.com/quant1x/exchange"
	"gitee.com/quant1x/gox/api"
	"sync"
)

var (
	wideKLinesMutex        sync.RWMutex
	routineLocalWideKLines = map[string][]SecurityFeature{}
)

// loadWideKLines 加载基础K线
func loadWideKLines(securityCode string) []SecurityFeature {
	filename := cache.WideFilename(securityCode)
	var lines []SecurityFeature
	_ = api.CsvToSlices(filename, &lines)
	return lines
}

func checkWideKLinesOffset(lines []SecurityFeature, date string) (offset int) {
	rows := len(lines)
	offset = 0
	for i := 0; i < rows; i++ {
		klineDate := lines[rows-1-i].Date
		if klineDate < date {
			return -1
		} else if klineDate == date {
			break
		} else {
			offset++
		}
	}
	if offset+1 >= rows {
		return -1
	}
	return
}

// updateCacheWideKLines 更新缓存宽表
func updateCacheWideKLines(securityCode string, lines []SecurityFeature) {
	if len(lines) == 0 {
		return
	}
	wideKLinesMutex.Lock()
	defer wideKLinesMutex.Unlock()
	routineLocalWideKLines[securityCode] = lines

}

// CheckoutWideKLines 捡出指定日期的K线数据
func CheckoutWideKLines(code, date string) []SecurityFeature {
	securityCode := exchange.CorrectSecurityCode(code)
	date = exchange.FixTradeDate(date)
	// 1. 取缓存的K线
	wideKLinesMutex.RLock()
	cacheKLines, ok := routineLocalWideKLines[securityCode]
	wideKLinesMutex.RUnlock()
	if !ok {
		cacheKLines = loadWideKLines(securityCode)
		updateCacheWideKLines(securityCode, cacheKLines)
	}
	rows := len(cacheKLines)
	if rows == 0 {
		return nil
	}
	// 1.1 检查是否最新数据
	kline := cacheKLines[rows-1]
	if kline.Date < date {
		// 数据太旧, 重新加载
		cacheKLines = loadWideKLines(securityCode)
		updateCacheWideKLines(securityCode, cacheKLines)
	}
	// 2. 对齐数据缓存的日期, 过滤可能存在停牌没有数据的情况
	offset := checkWideKLinesOffset(cacheKLines, date)
	if offset < 0 {
		return nil
	}
	// 3. 返回指定日期前的K线数据
	klines := cacheKLines[0 : rows-offset]
	return klines
}
