package datasets

import (
	"context"
	"gitee.com/quant1x/engine/cache"
	"gitee.com/quant1x/engine/datasets/base"
	"gitee.com/quant1x/gotdx"
	"gitee.com/quant1x/gotdx/proto"
	"gitee.com/quant1x/gotdx/quotes"
	"gitee.com/quant1x/gotdx/trading"
	"gitee.com/quant1x/gox/api"
	"gitee.com/quant1x/gox/logger"
	"gitee.com/quant1x/pandas"
	. "gitee.com/quant1x/pandas/formula"
	"gitee.com/quant1x/pandas/stat"
	"reflect"
	"strconv"
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
	//base.UpdateMinutes(this.GetSecurityCode(), date)
	GetKLineAll(this.GetSecurityCode())
}

func (this *DataWide) Repair(date string) {
	this.Update(date)
}

func (this *DataWide) Increase(snapshot quotes.Snapshot) {
	//TODO implement me
	panic("implement me")
}

func (this *DataWide) Print(code string, date ...string) {
	_ = code
	_ = date
}

var (
	FBarsProtocolFields = []string{"Open", "Close", "High", "Low", "Vol", "Amount", "DateTime", "UpCount", "DownCount"}
	FBarsRawFields      = []string{"open", "close", "high", "low", "volume", "amount", "date", "up", "down"}
	FBarsHalfFields     = []string{"date", "open", "close", "high", "low", "volume", "amount", "up", "down", "open_volume", "open_turnz", "open_unmatched", "close_volume", "close_turnz", "close_unmatched", "inner_volume", "outer_volume", "inner_amount", "outer_amount"}
	//FBarsOldWideFields  = []string{"date", "open", "close", "high", "low", "volume", "amount", "up", "down", "last_close", "turnover_rate", "open_volume", "open_turnz", "open_unmatched", "close_volume", "close_turnz", "close_unmatched", "inner_volume", "outer_volume", "inner_amount", "outer_amount"}
	FBarsWideFields = []string{"date", "open", "close", "high", "low", "volume", "amount", "up", "down", "last_close", "change_rate", "open_volume", "open_turnz", "open_unmatched", "close_volume", "close_turnz", "close_unmatched", "inner_volume", "outer_volume", "inner_amount", "outer_amount"}
)

//var (
//	// DataDaysDiff 日期差异偏移量
//	DataDaysDiff = 1
//)

// SetKLineOffset 设置K线数据调整回补天数
func SetKLineOffset(days int) {
	if days <= 1 {
		return
	}
	base.DataDaysDiff = days
}

// loadCacheKLine 加载K线
//
//	第2个参数, 是否前复权
func loadCacheKLine(code string, adjust ...bool) pandas.DataFrame {
	// 默认不复权
	qfq := false
	if len(adjust) > 0 {
		qfq = adjust[0]
	}
	filename := cache.FeatureFilename(code)
	var df pandas.DataFrame
	if !api.FileExist(filename) {
		return df
	} else {
		df = pandas.ReadCSV(filename)
	}
	// 调整字段流程
	{
		// turnover_rate 改为 change_rate
		df.SetName("turnover_rate", "change_rate")
	}

	fields := FBarsWideFields
	df = df.Select(fields)
	if df.Nrow() == 0 {
		return df
	}
	if qfq {
		xdxrs := base.GetCacheXdxrList(code)
		for i := 0; i < len(xdxrs); i++ {
			xdxr := xdxrs[i]
			if xdxr.Category != 1 {
				continue
			}
			xdxrDate := xdxr.Date
			factor := xdxr.Adjust()
			for j := 0; j < df.Nrow(); j++ {
				m1 := df.IndexOf(j, true)
				dt := m1["date"].(reflect.Value).String()
				if dt > xdxrDate {
					break
				}
				if dt < xdxrDate {
					po := m1["open"].(reflect.Value)
					po.SetFloat(factor(po.Float()))
					pc := m1["close"].(reflect.Value)
					pc.SetFloat(factor(pc.Float()))
					ph := m1["high"].(reflect.Value)
					ph.SetFloat(factor(ph.Float()))
					pl := m1["low"].(reflect.Value)
					pl.SetFloat(factor(pl.Float()))
				}
				plc := m1["last_close"].(reflect.Value)
				plc.SetFloat(factor(plc.Float()))
				if dt == xdxrDate {
					break
				}
			}
		}
	}
	return df
}

// GetCacheKLine 加载K线
//
//	第2个参数, 是否前复权
func GetCacheKLine(code string, adjust ...bool) pandas.DataFrame {
	df := loadCacheKLine(code, adjust...)
	if df.Nrow() == 0 {
		return df
	}
	// 取出成交量序列
	VOL := df.Col("volume")
	DATES := df.Col("date")
	lastDay := DATES.IndexOf(-1).(string)
	total := df.Nrow()
	// 计算5日均量
	mv5 := MA(VOL, 5)
	mav := REF(mv5, 1)
	lb := VOL.Div(mav)
	lb = lb.Apply2(func(idx int, v any) any {
		if idx+1 < total {
			return v
		} else {
			tmp := stat.Any2DType(v)
			ms := stat.DType(trading.Minutes(lastDay)) / float64(trading.CN_TOTALFZNUM)
			tmp /= ms
			return tmp
		}
	}, true)

	// 链接量比序列
	oLB := pandas.NewSeries(stat.SERIES_TYPE_DTYPE, "lb", lb.DTypes())
	oMV5 := pandas.NewSeries(stat.SERIES_TYPE_DTYPE, "mv5", mv5.Div(trading.CN_DEFAULT_TOTALFZNUM).DTypes())
	vr := VOL.Div(REF(VOL, 1))
	oVR := pandas.NewSeries(stat.SERIES_TYPE_DTYPE, "vr", vr.DTypes())
	CLOSE := df.Col("close")
	chg5 := CLOSE.Div(REF(CLOSE, 5)).Sub(1.00).Mul(100)
	chg10 := CLOSE.Div(REF(CLOSE, 10)).Sub(1.00).Mul(100)
	oChg5 := pandas.NewSeries(stat.SERIES_TYPE_FLOAT64, "chg5", chg5.DTypes())
	oChg10 := pandas.NewSeries(stat.SERIES_TYPE_FLOAT64, "chg10", chg10.DTypes())
	ma5 := MA(CLOSE, 5)
	ma10 := MA(CLOSE, 10)
	ma20 := MA(CLOSE, 20)
	oMA5 := pandas.NewSeries(stat.SERIES_TYPE_FLOAT64, "ma5", ma5.DTypes())
	oMA10 := pandas.NewSeries(stat.SERIES_TYPE_FLOAT64, "ma10", ma10.DTypes())
	oMA20 := pandas.NewSeries(stat.SERIES_TYPE_FLOAT64, "ma20", ma20.DTypes())
	AMOUNT := df.Col("amount")
	averagePrice := AMOUNT.Div(VOL)
	oAP := pandas.NewSeries(stat.SERIES_TYPE_FLOAT64, "ap", averagePrice.DTypes())
	df = df.Join(oLB, oMV5, oVR, oChg5, oChg10, oMA5, oMA10, oMA20, oAP)
	return df
}

// GetKLineAll 获取日K线
func GetKLineAll(securityCode string, argv ...int) pandas.DataFrame {
	kType := uint16(proto.KLINE_TYPE_RI_K)
	if len(argv) == 1 {
		kType = uint16(argv[0])
	}
	securityCode = proto.CorrectSecurityCode(securityCode)
	tdxApi := gotdx.GetTdxApi()
	startDate := proto.MARKET_CH_FIRST_LISTTIME
	// 默认是缓存中是已经复权的数据
	dfCache := loadCacheKLine(securityCode)
	isIndex := proto.AssertIndexBySecurityCode(securityCode)
	rawFields := FBarsProtocolFields
	newFields := FBarsRawFields
	// 尝试选择一次字段, 如果出现异常, 则清空dataframe, 重新下载
	dfCache = dfCache.Select(FBarsWideFields)
	if dfCache.Nrow() == 0 {
		dfCache = pandas.DataFrame{}
	}
	var info *quotes.FinanceInfo
	var err error
	var klineDaysOffset = base.DataDaysDiff
	if dfCache.Nrow() > 0 {
		ds := dfCache.Col("date").Strings()
		if klineDaysOffset > len(ds) {
			klineDaysOffset = len(ds)
		}
		startDate = ds[len(ds)-klineDaysOffset]
	} else {
		info, err = tdxApi.GetFinanceInfo(securityCode)
		if err != nil {
			return dfCache
		}
		if info.IPODate > 0 {
			startDate = strconv.FormatInt(int64(info.IPODate), 10)
			startDate = trading.FixTradeDate(startDate)
		}
	}
	endDate := trading.Today()
	ts := trading.TradeRange(startDate, endDate)
	history := []quotes.SecurityBar{}
	step := uint16(quotes.TDX_SECURITY_BARS_MAX)
	total := uint16(len(ts))
	start := uint16(0)
	hs := []quotes.SecurityBarsReply{}
	for {
		count := step
		if total-start >= step {
			count = step
		} else {
			count = total - start
		}
		var data *quotes.SecurityBarsReply
		var err error
		retryTimes := 0
		for retryTimes < quotes.DefaultRetryTimes {
			if isIndex {
				data, err = tdxApi.GetIndexBars(securityCode, kType, start, count)
			} else {
				data, err = tdxApi.GetKLine(securityCode, kType, start, count)
			}
			if err == nil && data != nil {
				break
			}
			retryTimes++
		}
		if err != nil {
			logger.Errorf("code=%s, error=%s", securityCode, err.Error())
			return pandas.DataFrame{}
		}
		hs = append(hs, *data)
		if data.Count < count {
			// 已经是最早的记录
			// 需要排序
			break
		}
		start += count
		if start >= total {
			break
		}
	}
	hs = stat.Reverse(hs)
	startDate = trading.FixTradeDate(startDate)
	for _, v := range hs {
		for _, row := range v.List {
			dateTime := trading.FixTradeDate(row.DateTime)
			if dateTime < startDate {
				continue
			}
			row.Vol = row.Vol * 100
			history = append(history, row)
		}
	}
	dfNew := pandas.LoadStructs(history)
	dfNew = dfNew.Select(rawFields)
	err = dfNew.SetNames(newFields...)
	if err != nil {
		return pandas.DataFrame{}
	}
	ds1 := dfNew.Col("date", true)
	ds1.Apply2(func(idx int, v any) any {
		date1 := v.(string)
		dt, err := api.ParseTime(date1)
		if err != nil {
			return date1
		}
		return dt.Format(trading.TradingDayDateFormat)
	}, true)
	// 补充昨日收盘和涨跌幅
	dfNew = attachVolume(dfNew, securityCode)
	dfNew = dfNew.Select(FBarsHalfFields)
	if dfNew.Nrow() > 0 {
		// 除权除息
		xdxrs := base.GetCacheXdxrList(securityCode)
		cacheLastDay := dfNew.Col("date").IndexOf(-1).(string)
		for i := 0; i < len(xdxrs); i++ {
			xdxr := xdxrs[i]
			if xdxr.Category != 1 || xdxr.Date < startDate || xdxr.Date > cacheLastDay {
				continue
			}
			xdxrDate := xdxr.Date
			factor := xdxr.Adjust()
			for j := 0; j < dfNew.Nrow(); j++ {
				m1 := dfNew.IndexOf(j, true)
				barCurrentDate := m1["date"].(reflect.Value).String()
				if barCurrentDate > xdxrDate {
					break
				}
				if barCurrentDate < xdxrDate {
					po := m1["open"].(reflect.Value)
					po.SetFloat(factor(po.Float()))
					pc := m1["close"].(reflect.Value)
					pc.SetFloat(factor(pc.Float()))
					ph := m1["high"].(reflect.Value)
					ph.SetFloat(factor(ph.Float()))
					pl := m1["low"].(reflect.Value)
					pl.SetFloat(factor(pl.Float()))
				}
				if barCurrentDate == xdxrDate {
					break
				}
			}
		}
	}
	dfCache = dfCache.Select(FBarsHalfFields)
	df := dfCache.Subset(0, dfCache.Nrow()-klineDaysOffset)
	if df.Nrow() > 0 {
		df = df.Concat(dfNew)
	} else {
		df = dfNew
	}
	CLOSE := df.Col("close")
	LAST := CLOSE.Shift(1)
	rate := CLOSE.Sub(LAST).Div(LAST).Mul(100.00).DTypes()

	lc := pandas.NewSeries(stat.SERIES_TYPE_FLOAT64, "last_close", LAST.DTypes())
	// TODO: turnover_rate 英文是换手率, 这里实际上是涨幅
	tr := pandas.NewSeries(stat.SERIES_TYPE_FLOAT64, "change_rate", rate)
	df = df.Join(lc, tr)
	df = df.Select(FBarsWideFields)
	if df.Nrow() > 0 {
		filename := cache.FeatureFilename(securityCode)
		_ = df.WriteCSV(filename)
	}
	return df
}

// 附加成交量
func attachVolume(df pandas.DataFrame, securityCode string) pandas.DataFrame {
	securityCode = proto.CorrectSecurityCode(securityCode)
	dates := df.Col("date").Strings()
	if len(dates) == 0 {
		return df
	}
	buyVolumes := []int64{}
	sellVolumes := []int64{}
	buyAmounts := []stat.DType{}
	sellAmounts := []stat.DType{}

	openVolumes := []int64{}
	openTurnZ := []stat.DType{}
	openUnmatched := []int64{}
	closeVolumes := []int64{}
	closeTurnZ := []stat.DType{}
	closeUnmatched := []int64{}
	for _, tradeDate := range dates {
		tmp := base.CheckoutTickData(securityCode, tradeDate, true)
		logger.Warnf("tick: code=%s, date=%s, rows=%d", securityCode, tradeDate, len(tmp))
		//summary := InflowCount(tmp, securityCode)
		summary := CountInflow(tmp, securityCode, tradeDate)
		buyVolumes = append(buyVolumes, summary.OuterVolume)
		sellVolumes = append(sellVolumes, summary.InnerVolume)
		buyAmounts = append(buyAmounts, summary.OuterAmount)
		sellAmounts = append(sellAmounts, summary.InnerAmount)

		openVolumes = append(openVolumes, summary.OpenVolume)
		openTurnZ = append(openTurnZ, summary.OpenTurnZ)
		openUnmatched = append(openUnmatched, summary.OpenUnmatched)
		closeVolumes = append(closeVolumes, summary.CloseVolume)
		closeTurnZ = append(closeTurnZ, summary.CloseTurnZ)
		closeUnmatched = append(closeUnmatched, summary.CloseUnmatched)
	}
	// 调整字段名
	bv := pandas.NewSeries(stat.SERIES_TYPE_INT64, "outer_volume", buyVolumes)
	sv := pandas.NewSeries(stat.SERIES_TYPE_INT64, "inner_volume", sellVolumes)
	ba := pandas.NewSeries(stat.SERIES_TYPE_DTYPE, "outer_amount", buyAmounts)
	sa := pandas.NewSeries(stat.SERIES_TYPE_DTYPE, "inner_amount", sellAmounts)

	// 新增字段
	ov := pandas.NewSeries(stat.SERIES_TYPE_INT64, "open_volume", openVolumes)
	ot := pandas.NewSeries(stat.SERIES_TYPE_FLOAT64, "open_turnz", openTurnZ)
	ou := pandas.NewSeries(stat.SERIES_TYPE_INT64, "open_unmatched", openUnmatched)
	cv := pandas.NewSeries(stat.SERIES_TYPE_INT64, "close_volume", closeVolumes)
	ct := pandas.NewSeries(stat.SERIES_TYPE_FLOAT64, "close_turnz", closeTurnZ)
	cu := pandas.NewSeries(stat.SERIES_TYPE_INT64, "close_unmatched", closeUnmatched)

	df = df.Join(bv, sv, ba, sa, ov, ot, ou, cv, ct, cu)
	return df
}