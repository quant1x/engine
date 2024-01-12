package factors

import (
	"context"
	"gitee.com/quant1x/engine/cache"
	"gitee.com/quant1x/engine/datasource/base"
	"gitee.com/quant1x/exchange"
	"gitee.com/quant1x/gotdx/quotes"
	"gitee.com/quant1x/gox/api"
	"slices"
)

type DataWide struct {
	Manifest
}

func init() {
	summary := mapDataSets[BaseKLineWide]
	_ = cache.Register(&DataWide{Manifest: Manifest{DataSummary: summary}})
}

func (this *DataWide) Clone(date string, code string) DataSet {
	summary := mapDataSets[BaseKLineWide]
	var dest = DataWide{
		Manifest: Manifest{
			DataSummary: summary,
			Date:        date,
			Code:        code,
		},
	}
	return &dest
}

func (this *DataWide) Init(ctx context.Context, date string) error {
	_ = ctx
	_ = date
	return nil
}

func (this *DataWide) Update(date string) {
	GetKLineAll(this.GetSecurityCode())
	_ = date
}

func (this *DataWide) Repair(date string) {
	this.Update(date)
}

func (this *DataWide) Increase(snapshot quotes.Snapshot) {
	_ = snapshot
}

func (this *DataWide) Print(code string, date ...string) {
	_ = code
	_ = date
}

// 通过日期拉取宽表数据
func pullWideByDate(securityCode, date string) []SecurityFeature {
	securityCode = exchange.CorrectSecurityCode(securityCode)
	// 1. 加载缓存
	filename := cache.WideFilename(securityCode)
	var list []SecurityFeature
	var beginDate string
	err := api.CsvToSlices(filename, &list)
	if err != nil || len(list) == 0 {
		beginDate = exchange.MARKET_CH_FIRST_LISTTIME
	} else {
		beginDate = list[len(list)-1].Date
	}
	// 2. 确定补齐数据的日期
	endDate := exchange.FixTradeDate(date)
	ts := exchange.TradingDateRange(beginDate, endDate)
	tslen := len(ts)
	if tslen == 0 {
		beginDate = endDate
	} else if tslen == 1 {
		beginDate = ts[0]
		endDate = ts[0]
	} else {
		beginDate = ts[0]
		endDate = ts[tslen-1]
	}
	// 3. 补齐数据
	list_length := len(list)
	// 3.1 补齐日线, 日线是必须要有的, 也肯定会有
	// 数据为空, 从基础K线获取K线部分
	klines := base.CheckoutKLines(securityCode, endDate)
	kline_length := len(klines)
	if kline_length == 0 {
		// K线为空, 返回空
		return nil
	} else {
		//beginDate = klines[0].Date
	}
	// 如果kline比wide数据多
	if kline_length > list_length {
		list = slices.Grow(list, kline_length)
	}
	transBeginDate := base.GetTickStartDate()
	transBeginDate = exchange.FixTradeDate(transBeginDate)
	for i, v := range klines {
		if v.Date < beginDate {
			continue
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
	// 4 保存文件
	_ = api.SlicesToCsv(filename, list)
	return list
}
