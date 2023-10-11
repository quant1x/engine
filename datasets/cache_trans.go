package datasets

import (
	"gitee.com/quant1x/engine/cache"
	"gitee.com/quant1x/gotdx/quotes"
)

// TransactionRecord 成交记录
//
//	最短3秒内的合并统计数据, 与行情数据保持一致
//	不可以当作tick数据来使用
type TransactionRecord struct {
	DataCache
}

func init() {
	_ = cache.Register(&TransactionRecord{})
}

func (r *TransactionRecord) Get(code string, date ...string) any {
	//TODO implement me
	panic("implement me")
}

func (r *TransactionRecord) Kind() DataKind {
	return BaseTransaction
}

func (r *TransactionRecord) Name() string {
	return mapDataSets[r.Kind()].Name
}

func (r *TransactionRecord) Key() string {
	return mapDataSets[r.Kind()].Key
}

func (r *TransactionRecord) Init(barIndex *int, date string) error {
	_ = barIndex
	_ = date
	return nil
}

func (r *TransactionRecord) Filename(date, code string) string {
	//TODO implement me
	panic("implement me")
}

func (r *TransactionRecord) Update(cacheDate, featureDate string) {
	//TODO implement me
	panic("implement me")
}

func (r *TransactionRecord) Repair(cacheDate, featureDate string) {
	//TODO implement me
	panic("implement me")
}

func (r *TransactionRecord) Increase(snapshot quotes.Snapshot) {
	//TODO implement me
	panic("implement me")
}

func (r *TransactionRecord) Clone(date string, code string) DataSet {
	var dest = TransactionRecord{DataCache: DataCache{Date: date, Code: code}}
	return &dest
}

//// GetTickAll 下载全部tick数据
//func GetTickAll(securityCode string) {
//	defer func() {
//		// 解析失败以后输出日志, 以备检查
//		if err := recover(); err != nil {
//			logger.Errorf("下载tick数据异常: code=%s", securityCode)
//			return
//		}
//	}()
//	securityCode = proto.CorrectSecurityCode(securityCode)
//	tdxApi := gotdx.GetTdxApi()
//	info, err := tdxApi.GetFinanceInfo(securityCode)
//	if err != nil {
//		return
//	}
//	tStart := strconv.FormatInt(int64(info.IPODate), 10)
//	fixStart := __tickHistoryStartDate
//	if tStart < fixStart {
//		tStart = fixStart
//	}
//	tEnd := trading.Today()
//	dateRange := trading.TradeRange(tStart, tEnd)
//	// 反转切片
//	dateRange = stat.Reverse(dateRange)
//	if len(dateRange) == 0 {
//		return
//	}
//	bar := progressbar.NewBar(2, fmt.Sprintf("同步[%s]", securityCode), len(dateRange))
//	today := dateRange[0]
//	ignore := false
//	for _, tradeDate := range dateRange {
//		bar.Add(1)
//		if ignore {
//			continue
//		}
//		fname := cache.TickFilename(securityCode, tradeDate)
//		if tradeDate != today && api.FileIsValid(fname) {
//			// 如果已经存在, 假定之前的数据已经下载过了, 不需要继续
//			ignore = true
//			continue
//		}
//		df := GetTickData(securityCode, tradeDate)
//		if df.Nrow() == 0 && tradeDate != today {
//			// 如果数据为空, 且不是当前日期, 认定为从这天起往前是没有分笔成交数据的
//			ignore = true
//		}
//	}
//
//	return
//}
