package base

import (
	"gitee.com/quant1x/engine/cache"
	"gitee.com/quant1x/exchange"
	"gitee.com/quant1x/gotdx"
	"gitee.com/quant1x/gotdx/quotes"
	"gitee.com/quant1x/gox/api"
	"gitee.com/quant1x/gox/logger"
	"gitee.com/quant1x/pandas/stat"
	"strconv"
)

const (
	TradingFirstTime        = "09:25"      // 第一个时间
	TradingStartTime        = "09:30"      // 开盘时间
	TradingFinalBiddingTime = "14:57"      // 尾盘集合竞价时间
	TradingLastTime         = "15:00"      // 最后一个时间
	TickDefaultStartDate    = "2023-01-01" // 分笔成交最早的日期
)

var (
	// TickDefaultStartDate 最早的时间
	__tickHistoryStartDate = "20230101"
)

// UpdateTickStartDate 修改tick数据开始下载的日期
func UpdateTickStartDate(date string) {
	dt, err := api.ParseTime(date)
	if err != nil {
		return
	}
	date = dt.Format(cache.TDX_FORMAT_PROTOCOL_DATE)
	__tickHistoryStartDate = date
}

func GetTickStartDate() string {
	return __tickHistoryStartDate
}

// Transaction 获取指定日期的历史成交数据
func Transaction(securityCode, tradeDate string) []quotes.TickTransaction {
	securityCode = exchange.CorrectSecurityCode(securityCode)
	tdxApi := gotdx.GetTdxApi()
	offset := uint16(quotes.TDX_TRANSACTION_MAX)
	start := uint16(0)
	history := make([]quotes.TickTransaction, 0)
	hs := make([]quotes.TransactionReply, 0)
	date := exchange.FixTradeDate(tradeDate, cache.TDX_FORMAT_PROTOCOL_DATE)
	iDate := stat.AnyToInt64(date)
	for {
		var data *quotes.TransactionReply
		var err error
		retryTimes := 0
		for retryTimes < quotes.DefaultRetryTimes {
			data, err = tdxApi.GetHistoryTransactionData(securityCode, uint32(iDate), start, offset)
			if err == nil && data != nil {
				break
			}
			retryTimes++
		}
		if err != nil {
			logger.Errorf("code=%s, tradeDate=%s, error=%s", securityCode, tradeDate, err.Error())
			return []quotes.TickTransaction{}
		}
		if data == nil || data.Count == 0 {
			break
		}
		// 历史成交记录是按照时间排序
		//data.List = stat.Reverse(data.List)
		hs = append(hs, *data)
		if data.Count < offset {
			break
		}
		start += offset
	}
	// 这里需要反转一下
	hs = stat.Reverse(hs)
	for _, v := range hs {
		history = append(history, v.List...)
	}

	return history
}

// GetTickAll 下载全部tick数据
func GetTickAll(securityCode string) {
	defer func() {
		// 解析失败以后输出日志, 以备检查
		if err := recover(); err != nil {
			logger.Errorf("下载tick数据异常: code=%s", securityCode)
			return
		}
	}()
	securityCode = exchange.CorrectSecurityCode(securityCode)
	tdxApi := gotdx.GetTdxApi()
	info, err := tdxApi.GetFinanceInfo(securityCode)
	if err != nil {
		return
	}
	tStart := strconv.FormatInt(int64(info.IPODate), 10)
	fixStart := __tickHistoryStartDate
	//fmt.Println("start date:", fixStart)
	if tStart < fixStart {
		tStart = fixStart
	}
	tEnd := exchange.Today()
	dateRange := exchange.TradeRange(tStart, tEnd)
	// 反转切片
	dateRange = stat.Reverse(dateRange)
	if len(dateRange) == 0 {
		return
	}
	today := dateRange[0]
	ignore := false
	for _, tradeDate := range dateRange {
		if ignore {
			continue
		}
		fname := cache.TickFilename(securityCode, tradeDate)
		if tradeDate != today && api.FileIsValid(fname) {
			// 如果已经存在, 假定之前的数据已经下载过了, 不需要继续
			ignore = true
			continue
		}
		list := GetTickData(securityCode, tradeDate)
		if len(list) == 0 && tradeDate != today {
			// 如果数据为空, 且不是当前日期, 认定为从这天起往前是没有分笔成交数据的
			ignore = true
		}
	}

	return
}

// GetTickData 获取指定日期的分笔成交记录
func GetTickData(securityCode string, date string) (list []quotes.TickTransaction) {
	securityCode = exchange.CorrectSecurityCode(securityCode)
	list = CheckoutTickData(securityCode, date, false)
	if len(list) == 0 {
		return list
	}
	tickFile := cache.TickFilename(securityCode, date)
	err := api.SlicesToCsv(tickFile, list)
	if err != nil {
		return []quotes.TickTransaction{}
	}

	return list
}

// CheckoutTickData 获取指定日期的分笔成交记录
//
//	先从缓存获取, 如果缓存不存在, 则从服务器下载
//	K线附加成交数据
func CheckoutTickData(securityCode string, date string, ignorePreviousData bool) (list []quotes.TickTransaction) {
	securityCode = exchange.CorrectSecurityCode(securityCode)
	// 对齐日期格式: YYYYMMDD
	tradeDate := exchange.FixTradeDate(date, cache.TDX_FORMAT_PROTOCOL_DATE)
	if ignorePreviousData {
		// 在默认日期之前的数据直接返回空
		startDate := exchange.FixTradeDate(__tickHistoryStartDate, cache.TDX_FORMAT_PROTOCOL_DATE)
		if tradeDate < startDate {
			logger.Errorf("tick: code=%s, trade-date=%s, start-date=%s, 没有数据", securityCode, tradeDate, startDate)
			return list
		}
	}
	startTime := TradingFirstTime
	filename := cache.TickFilename(securityCode, tradeDate)
	if api.FileExist(filename) {
		// 如果缓存存在
		err := api.CsvToSlices(filename, &list)
		cacheLength := len(list)
		if err == nil && cacheLength > 0 {
			lastTime := list[cacheLength-1].Time
			if lastTime == TradingLastTime {
				//logger.Warnf("tick: code=%s, trade-date=%s, 缓存存在", securityCode, tradeDate)
				return
			}
			firstTime := ""
			skipCount := 0
			for i := 0; i < cacheLength; i++ {
				tm := list[cacheLength-1-i].Time
				if len(firstTime) == 0 {
					firstTime = tm
					startTime = firstTime
					skipCount++
					continue
				}
				if tm < firstTime {
					startTime = firstTime
					break
				} else {
					skipCount++
				}
			}
			// 截取 startTime之前的记录
			list = list[0 : cacheLength-skipCount]
		} else {
			logger.Errorf("tick: code=%s, trade-date=%s, 没有有效数据, %+v", securityCode, tradeDate, err)
		}
	}

	tdxApi := gotdx.GetTdxApi()
	offset := uint16(quotes.TDX_TRANSACTION_MAX)
	// 只求增量, 分笔成交数据是从后往前取数据, 缓存是从前到后顺序存取
	start := uint16(0)
	history := make([]quotes.TickTransaction, 0)
	hs := make([]quotes.TransactionReply, 0)
	for {
		var data *quotes.TransactionReply
		var err error
		retryTimes := 0
		for retryTimes < quotes.DefaultRetryTimes {
			if exchange.CurrentlyTrading(tradeDate) {
				data, err = tdxApi.GetTransactionData(securityCode, start, offset)
			} else {
				data, err = tdxApi.GetHistoryTransactionData(securityCode, toTdxProtocolDate(tradeDate), start, offset)
			}
			if err == nil && data != nil {
				break
			}
			retryTimes++
		}
		if err != nil {
			logger.Errorf("code=%s, tradeDate=%s, error=%s", securityCode, tradeDate, err.Error())
			return
		}
		if data == nil || data.Count == 0 {
			break
		}
		var tmp quotes.TransactionReply
		tmpList := stat.Reverse(data.List)
		for _, td := range tmpList {
			// 追加包含startTime之后的记录
			if td.Time >= startTime {
				tmp.Count += 1
				tmp.List = append(tmp.List, td)
			}
		}
		tmp.List = stat.Reverse(tmp.List)
		hs = append(hs, tmp)
		if tmp.Count < offset {
			// 已经是最早的记录
			// 需要排序
			break
		}
		start += offset
	}
	hs = stat.Reverse(hs)
	for _, v := range hs {
		history = append(history, v.List...)
	}
	if len(history) == 0 {
		return
	}
	list = append(list, history...)

	return
}
