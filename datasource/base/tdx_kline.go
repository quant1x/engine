package base

import (
	"strings"
	"sync"
	"time"

	"github.com/quant1x/exchange"
	"github.com/quant1x/num"
	"github.com/quant1x/pandas"
	"github.com/quant1x/x/api"
)

var (
	klineMutex         sync.RWMutex
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

// UpdateCacheKLines 更新缓存K线
func UpdateCacheKLines(securityCode string, klines []KLine) {
	if len(klines) == 0 {
		return
	}
	klineMutex.Lock()
	routineLocalKLines[securityCode] = klines
	klineMutex.Unlock()
}

// CheckoutKLines 捡出指定日期的K线数据
func CheckoutKLines(code, date string) []KLine {
	securityCode := exchange.CorrectSecurityCode(code)
	date = exchange.FixTradeDate(date)
	// 1. 取缓存的K线
	klineMutex.RLock()
	cacheKLines, ok := routineLocalKLines[securityCode]
	klineMutex.RUnlock()
	if !ok {
		cacheKLines = LoadBasicKline(securityCode)
		UpdateCacheKLines(securityCode, cacheKLines)
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
		UpdateCacheKLines(securityCode, cacheKLines)
	}
	// 2. 对齐数据缓存的日期, 过滤可能存在停牌没有数据的情况
	offset := checkKLineOffset(cacheKLines, date)
	if offset < 0 {
		return nil
	}
	// 3. 返回指定日期前的K线数据
	klines := cacheKLines[0 : rows-offset]
	return klines
}

// 由日线计算日线以上级别的K线
func getPeriodKLine(checkPeriod func(date ...string) (s, e string), securityCode string, cacheKLine ...[]KLine) (list []KLine) {
	baseKLines := []KLine{}
	if len(cacheKLine) > 0 {
		baseKLines = cacheKLine[0]
	} else {
		baseKLines = LoadBasicKline(securityCode)
	}
	if len(baseKLines) == 0 {
		return
	}

	var klines []KLine
	var kline KLine
	length := len(baseKLines)
	for i, v := range baseKLines {
		// 确定时间, 周线的日期是本周内最后一个交易日
		if len(kline.Date) == 0 {
			// 重新计算周线, 先确认周线范围
			ws, we := checkPeriod(v.Date)
			_ = ws
			periodLastDate := exchange.FixTradeDate(we)
			offset := i
			for {
				destDate := baseKLines[offset].Date
				if destDate < periodLastDate {
					offset++
				}
				if offset >= length {
					periodLastDate = destDate
					break
				} else if baseKLines[offset].Date == periodLastDate {
					periodLastDate = baseKLines[offset].Date
					break
				} else if baseKLines[offset].Date > periodLastDate {
					periodLastDate = destDate
					break
				}
			}
			kline.Date = periodLastDate
		}
		// 周线开盘价以第一天OPEN为准
		if kline.Open == num.DType(0) {
			kline.Open = v.Open
		}
		// 周线的收盘价以本周最后一个交易日的CLOSE为准
		kline.Close = v.Close
		if kline.High == num.DType(0) {
			kline.High = v.High
		} else if kline.High < v.High {
			kline.High = v.High
		}
		if kline.Low == num.DType(0) {
			kline.Low = v.Low
		} else if kline.Low > v.Low {
			kline.Low = v.Low
		}
		kline.Volume += v.Volume
		kline.Amount += v.Amount

		// 切换下一周
		if kline.Date == v.Date || i+1 >= len(baseKLines) {
			kline.Date = v.Date
			klines = append(klines, kline)
			kline = KLine{}
		}
	}
	return klines
}

// 由日线计算日线以上级别的K线
func periodKLine(checkPeriod func(date ...string) (s, e string), securityCode string, cacheKLine ...[]KLine) (df pandas.DataFrame) {
	klines := getPeriodKLine(checkPeriod, securityCode, cacheKLine...)
	df = pandas.LoadStructs(klines)
	return
}

// WKLine 周线
func WKLine(securityCode string, cacheKLine ...[]KLine) []KLine {
	klines := getPeriodKLine(api.GetWeekDay, securityCode, cacheKLine...)
	return klines
}

// WeeklyKLine 周线
func WeeklyKLine(securityCode string, cacheKLine ...[]KLine) (df pandas.DataFrame) {
	return periodKLine(api.GetWeekDay, securityCode, cacheKLine...)
}

// MonthlyKLine 月K线
func MonthlyKLine(securityCode string, cacheKLine ...[]KLine) (df pandas.DataFrame) {
	return periodKLine(api.GetMonthDay, securityCode, cacheKLine...)
}

// 通过给定的日期 获得日期所在周的初始和结束日期
func getWeekDay(date string) (first, last string) {
	return api.GetWeekDay(date)
}

// 通过给定的日期 获得日期所在月的初始和结束日期
func getMonthDay(date string) (first, last string) {
	return api.GetMonthDay(date)
}

// 通过给定的日期 获得日期所在季度的初始和结束日期
func getQuarterDay(date string) (first, last string) {
	diff := 0
	_, first, last = api.GetQuarterByDate(date, diff)
	return first, last
}

// 通过给定的日期 获得日期所在年的初始和结束日期
func getYearDay(date string) (first, last string) {
	now := time.Now()
	if len(date) > 0 {
		day, err := api.ParseTime(date)
		if err == nil {
			now = day
		}
	}
	currentYear, _, _ := now.Date()
	currentLocation := now.Location()

	firstOfYear := time.Date(currentYear, 1, 1, 0, 0, 0, 0, currentLocation)
	lastOfYear := firstOfYear.AddDate(1, 0, -1)
	f := firstOfYear.Unix()
	l := lastOfYear.Unix()
	return time.Unix(f, 0).Format("2006-01-02") + " 00:00:00", time.Unix(l, 0).Format("2006-01-02") + " 23:59:59"
}

// ConvertKlinesTrading 转换K线周期
func ConvertKlinesTrading(cacheKLine []KLine, period string) []KLine {
	period = strings.ToLower(period)
	var checkPeriod func(date string) (s, e string)
	switch period {
	case "w":
		checkPeriod = getWeekDay
	case "m":
		checkPeriod = getMonthDay
	case "q":
		checkPeriod = getQuarterDay
	case "y":
		checkPeriod = getYearDay
	default:
		return cacheKLine
	}

	var klines []KLine
	var kline KLine
	length := len(cacheKLine)
	for i, v := range cacheKLine {
		// 确定时间, 周线的日期是本周内最后一个交易日
		if len(kline.Date) == 0 {
			// 重新计算周线, 先确认周线范围
			ws, we := checkPeriod(v.Date)
			_ = ws
			periodLastDate := exchange.FixTradeDate(we)
			offset := i
			for {
				destDate := cacheKLine[offset].Date
				if destDate < periodLastDate {
					offset++
				}
				if offset >= length {
					periodLastDate = destDate
					break
				} else if cacheKLine[offset].Date == periodLastDate {
					periodLastDate = cacheKLine[offset].Date
					break
				} else if cacheKLine[offset].Date > periodLastDate {
					periodLastDate = destDate
					break
				}
			}
			kline.Date = periodLastDate
		}
		// 周线开盘价以第一天OPEN为准
		if kline.Open == num.DType(0) {
			kline.Open = v.Open
		}
		// 周线的收盘价以本周最后一个交易日的CLOSE为准
		kline.Close = v.Close
		if kline.High == num.DType(0) {
			kline.High = v.High
		} else if kline.High < v.High {
			kline.High = v.High
		}
		if kline.Low == num.DType(0) {
			kline.Low = v.Low
		} else if kline.Low > v.Low {
			kline.Low = v.Low
		}
		kline.Volume += v.Volume
		kline.Amount += v.Amount

		// 切换下一周
		if kline.Date == v.Date || i+1 >= len(cacheKLine) {
			kline.Date = v.Date
			klines = append(klines, kline)
			kline = KLine{}
		}
	}
	return klines
}
