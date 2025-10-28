package factors

import (
	"gitee.com/quant1x/data/exchange"
	"gitee.com/quant1x/engine/cache"
	"gitee.com/quant1x/gox/api"
)

//var (
//	wideTableMutex        sync.RWMutex
//	routineLocalWideTable = map[string][]SecurityFeature{}
//)
//
//// updateCacheWideKLines 更新缓存宽表
//func updateCacheWideKLines(securityCode string, lines []SecurityFeature) {
//	if len(lines) == 0 {
//		return
//	}
//	wideTableMutex.Lock()
//	defer wideTableMutex.Unlock()
//	routineLocalWideTable[securityCode] = lines
//
//}
//
//// CheckoutWideKLines 捡出指定日期的K线数据
//func CheckoutWideKLines(code, date string) []SecurityFeature {
//	securityCode := exchange.CorrectSecurityCode(code)
//	date = exchange.FixTradeDate(date)
//	// 1. 取缓存的K线
//	wideTableMutex.RLock()
//	cacheLines, ok := routineLocalWideTable[securityCode]
//	wideTableMutex.RUnlock()
//	if !ok {
//		cacheLines = loadWideTable(securityCode)
//		updateCacheWideKLines(securityCode, cacheLines)
//	}
//	rows := len(cacheLines)
//	if rows == 0 {
//		return nil
//	}
//	// 1.1 检查是否最新数据
//	kline := cacheLines[rows-1]
//	if kline.Date < date {
//		// 数据太旧, 重新加载
//		cacheLines = loadWideTable(securityCode)
//		updateCacheWideKLines(securityCode, cacheLines)
//	}
//	// 2. 对齐数据缓存的日期, 过滤可能存在停牌没有数据的情况
//	offset := checkWideTableOffset(cacheLines, date)
//	if offset < 0 {
//		return nil
//	}
//	// 3. 返回指定日期前的K线数据
//	lines := cacheLines[0 : rows-offset]
//	return lines
//}

// loadWideTable 加载基础K线
func loadWideTable(securityCode string) []SecurityFeature {
	filename := cache.WideFilename(securityCode)
	var lines []SecurityFeature
	_ = api.CsvToSlices(filename, &lines)
	return lines
}

// 矫正日期的偏移
func checkWideTableOffset(lines []SecurityFeature, date string) (offset int) {
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

// CheckoutWideTableByDate 捡出指定日期的K线数据
//
//	TODO: 扩展数据不适用内存缓存, 减小内存压力
func CheckoutWideTableByDate(code, date string) []SecurityFeature {
	securityCode := exchange.CorrectSecurityCode(code)
	date = exchange.FixTradeDate(date)
	cacheLines := loadWideTable(securityCode)
	rows := len(cacheLines)
	if rows == 0 {
		return nil
	}
	// 2. 对齐数据缓存的日期, 过滤可能存在停牌没有数据的情况
	offset := checkWideTableOffset(cacheLines, date)
	if offset < 0 {
		return nil
	}
	// 3. 返回指定日期前的K线数据
	lines := cacheLines[0 : rows-offset]
	return lines
}
