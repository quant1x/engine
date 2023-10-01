package base

import (
	"gitee.com/quant1x/gotdx/proto"
	"gitee.com/quant1x/gotdx/trading"
)

var (
	// TODO: 隔日需要清空重新缓存
	routineLocalKLines = map[string][]KLine{}
)

func checkKLineOffset(klines []KLine, date string) (offset int) {
	rows := len(klines)
	offset = 0
	for i := 0; i < rows; i++ {
		klineDate := klines[rows-1-i].Date
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

func CheckoutKLines(code, date string) []KLine {
	securityCode := proto.CorrectSecurityCode(code)
	date = trading.FixTradeDate(date)
	// 1. 取缓存的K线
	cacheKLines, ok := routineLocalKLines[securityCode]
	if !ok {
		cacheKLines = LoadBasicKline(securityCode)
	}
	rows := len(cacheKLines)
	if rows == 0 {
		return nil
	}
	// 1.1 检查是否最新数据
	kline := cacheKLines[rows-1]
	if kline.Date < date {
		// 数据太旧, 重新加载
		cacheKLines = LoadBasicKline(securityCode)
	}
	// 1.2 覆盖缓存
	routineLocalKLines[securityCode] = cacheKLines
	// 2. 对齐数据缓存的日期, 过滤可能存在停牌没有数据的情况
	offset := checkKLineOffset(cacheKLines, date)
	if offset < 0 {
		return nil
	}
	// 3. 返回指定日期前的K线数据
	klines := cacheKLines[0 : rows-offset]
	return klines
}
