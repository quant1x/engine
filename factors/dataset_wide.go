package factors

import (
	"context"
	"gitee.com/quant1x/engine/cache"
	"gitee.com/quant1x/engine/datasource/base"
	"gitee.com/quant1x/exchange"
	"gitee.com/quant1x/gotdx/quotes"
	"gitee.com/quant1x/gox/api"
	"gitee.com/quant1x/gox/num"
	"slices"
)

// DataWideKLine 宽表
type DataWideKLine struct {
	Manifest
}

func init() {
	summary := mapDataSets[BaseWideKLine]
	_ = cache.Register(&DataWideKLine{Manifest: Manifest{DataSummary: summary}})
}

func (this *DataWideKLine) Clone(date string, code string) DataSet {
	summary := mapDataSets[BaseWideKLine]
	var dest = DataWideKLine{
		Manifest: Manifest{
			DataSummary: summary,
			Date:        date,
			Code:        code,
		},
	}
	return &dest
}

func (this *DataWideKLine) Init(ctx context.Context, date string) error {
	_ = ctx
	_ = date
	return nil
}

func (this *DataWideKLine) Update(date string) {
	pullWideByDate(this.GetSecurityCode(), date)
}

func (this *DataWideKLine) Repair(date string) {
	this.Update(date)
}

func (this *DataWideKLine) Increase(snapshot quotes.Snapshot) {
	_ = snapshot
}

func (this *DataWideKLine) Print(code string, date ...string) {
	_ = code
	_ = date
}

// 通过日期拉取宽表数据
func pullWideByDate(securityCode, date string) []SecurityFeature {
	securityCode = exchange.CorrectSecurityCode(securityCode)
	// 1. 加载缓存
	filename := cache.WideFilename(securityCode)
	var list []SecurityFeature
	var beginDate string // 补数据的开始日期
	var endDate string   // 补数据的结束日期
	var cacheBeginDate, cacheEndDate string
	err := api.CsvToSlices(filename, &list)
	if err != nil || len(list) == 0 {
		// 如果文件为空, 暂定从1990-12-19
		cacheBeginDate = exchange.MARKET_CH_FIRST_LISTTIME
		cacheEndDate = cacheBeginDate
		beginDate = cacheBeginDate
	} else {
		cacheBeginDate = list[0].Date
		last := list[len(list)-1]
		cacheEndDate = last.Date
		// 以缓存文件最后一条记录的日期
		beginDate = cacheEndDate

	}
	// 2. 确定补齐数据的日期
	endDate = exchange.FixTradeDate(date)
	// 2.1 结束日期经过交易日历的校对处理一次
	//logger.Warnf("[%s]: begin=%s, end= %s", securityCode, beginDate, endDate)
	if len(beginDate) == 0 {
		beginDate = exchange.MARKET_CH_FIRST_LISTTIME
	}
	dates := exchange.TradingDateRange(beginDate, endDate)
	n := len(dates)
	if n == 0 {
		// 这种情况的原因是传入的date小于缓存的最后一条记录的日期
		beginDate = endDate
	} else if n == 1 {
		// 传入日期和缓存最后一条记录的日期相同
		beginDate = dates[0]
		endDate = dates[0]
	} else {
		beginDate = dates[0]
		endDate = dates[n-1]
	}
	// 3. 补齐日线, 日线是必须要有的, 也肯定会有
	// 数据为空, 从基础K线获取K线部分
	klines := base.CheckoutKLines(securityCode, endDate)
	kline_length := len(klines)
	if kline_length == 0 {
		// K线为空, 返回空
		return nil
	} else {
		// 校验wide缓存和k线缓存的开始日期是否对齐
		klsBeginDate := klines[0].Date
		klsEndDate := klines[kline_length-1].Date
		if cacheBeginDate == klsBeginDate {
			// 如果缓存的开始日期和k线的开始日期相同, 没有问题
		} else {
			// 如果缓存的开始日期和k线的开始日期不同, 则认为数据错乱, 清空
			clear(list)
			// 设置缓存开始日期为k线的开始日期
			cacheBeginDate = klsBeginDate
			// 设置缓存结束日期为k线的开始日期
			cacheEndDate = klsBeginDate
			beginDate = cacheBeginDate
		}
		_ = klsEndDate
	}
	// 4. 确定缓存记录数
	list_length := len(list)
	// 5. 如果kline比wide数据多
	if kline_length > list_length {
		list = slices.Grow(list, kline_length)
	}
	transBeginDate := base.GetBeginDateOfHistoricalTradingData()
	transBeginDate = exchange.FixTradeDate(transBeginDate)
	for i, v := range klines {
		if v.Date < beginDate {
			continue
		}
		if v.Date > endDate {
			break
		}
		featureDate := v.Date
		cacheDate := v.Date
		var info SecurityFeature
		// 复制k线
		info.Date = v.Date
		info.Open = v.Open
		info.Close = v.Close
		info.High = v.High
		info.Low = v.Low
		info.Volume = int64(v.Volume)
		info.Amount = v.Amount
		info.Up = v.Up
		info.Down = v.Down
		// 附加成交数据
		if featureDate >= transBeginDate {
			// 成交数据
			trans := base.CheckoutTransactionData(securityCode, featureDate, true)
			if len(list) > 0 {
				cover := CountInflow(trans, securityCode, featureDate)
				// 修正f10的缓存, 应该是缓存日期为准
				f10 := GetL5F10(securityCode, cacheDate)
				if f10 != nil {
					cover.OpenTurnZ = f10.TurnZ(cover.OpenVolume)
					cover.CloseTurnZ = f10.TurnZ(cover.CloseVolume)
				}
				info.OpenVolume = cover.OpenVolume
				info.OpenTurnZ = cover.OpenTurnZ
				info.CloseVolume = cover.CloseVolume
				info.CloseTurnZ = cover.CloseTurnZ
				info.Volume = cover.InnerVolume + cover.OuterVolume
				info.InnerVolume = cover.InnerVolume
				info.OuterVolume = cover.OuterVolume
				info.InnerAmount = cover.InnerAmount
				info.OuterAmount = cover.OuterAmount
			}
		}
		if i < list_length {
			list[i] = info
		} else {
			list = append(list, info)
		}
	}
	// 6. 修正last_close和change_rate
	lastClose := 0.000
	for i := 0; i < len(list); i++ {
		v := &list[i]
		if i == 0 {
			v.LastClose = v.Open
		} else {
			v.LastClose = lastClose
		}
		v.ChangeRate = num.NetChangeRate(v.LastClose, v.Close)
		lastClose = v.Close
	}

	// 7. 保存文件
	_ = api.SlicesToCsv(filename, list)
	return list
}
