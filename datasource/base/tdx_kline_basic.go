package base

import (
	"gitee.com/quant1x/engine/cache"
	"gitee.com/quant1x/exchange"
	"gitee.com/quant1x/gotdx"
	"gitee.com/quant1x/gotdx/proto"
	"gitee.com/quant1x/gotdx/quotes"
	"gitee.com/quant1x/gox/api"
	"gitee.com/quant1x/gox/logger"
)

var (
	// DataDaysDiff 日期差异偏移量
	DataDaysDiff = 1
)

// KLine 日K线基础结构
type KLine struct {
	Date     string  `name:"日期" dataframe:"date"`       // 日期
	Open     float64 `name:"开盘" dataframe:"open"`       // 开盘价
	Close    float64 `name:"收盘" dataframe:"close"`      // 收盘价
	High     float64 `name:"最高" dataframe:"high"`       // 最高价
	Low      float64 `name:"最低" dataframe:"low"`        // 最低价
	Volume   float64 `name:"成交量(股)" dataframe:"volume"` // 成交量
	Amount   float64 `name:"成交额(元)" dataframe:"amount"` // 成交金额
	Up       int     `name:"上涨/外盘" dataframe:"up"`      // 上涨家数
	Down     int     `name:"下跌/内盘" dataframe:"down"`    // 下跌家数
	Datetime string  `name:"时间" dataframe:"datetime"`   // 时间
}

// LoadBasicKline 加载基础K线
func LoadBasicKline(securityCode string) []KLine {
	filename := cache.KLineFilename(securityCode)
	var klines []KLine
	_ = api.CsvToSlices(filename, &klines)
	return klines
}

// UpdateAllBasicKLine 更新全部日K线基础数据并保存文件
func UpdateAllBasicKLine(securityCode string) []KLine {
	startDate := exchange.MARKET_CN_FIRST_DATE
	securityCode = exchange.CorrectSecurityCode(securityCode)
	isIndex := exchange.AssertIndexBySecurityCode(securityCode)
	cacheKLines := LoadBasicKline(securityCode)
	kLength := len(cacheKLines)
	var klineDaysOffset = DataDaysDiff
	if kLength > 0 {
		if klineDaysOffset > kLength {
			klineDaysOffset = kLength
		}
		startDate = cacheKLines[kLength-klineDaysOffset].Date
	} else {
		//f10 := flash.GetL5F10(securityCode)
		//if f10 != nil && len(f10.IpoDate) > 0 {
		//	startDate = f10.IpoDate
		//	startDate = trading.FixTradeDate(startDate)
		//}
	}
	endDate := exchange.Today()
	ts := exchange.TradeRange(startDate, endDate)
	history := make([]quotes.SecurityBar, 0)
	step := uint16(quotes.TDX_SECURITY_BARS_MAX)
	total := uint16(len(ts))
	start := uint16(0)
	hs := make([]quotes.SecurityBarsReply, 0)
	kType := uint16(proto.KLINE_TYPE_RI_K)
	tdxApi := gotdx.GetTdxApi()
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
			return []KLine{}
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
	hs = api.Reverse(hs)
	startDate = exchange.FixTradeDate(startDate)
	for _, v := range hs {
		for _, row := range v.List {
			dateTime := exchange.FixTradeDate(row.DateTime)
			if dateTime < startDate {
				continue
			}
			row.Vol = row.Vol * 100
			history = append(history, row)
		}
	}
	var newKLines []KLine
	for _, v := range history {
		date := exchange.FixTradeDate(v.DateTime)
		kline := KLine{
			Date:     date,
			Open:     v.Open,
			Close:    v.Close,
			High:     v.High,
			Low:      v.Low,
			Volume:   v.Vol,
			Amount:   v.Amount,
			Up:       int(v.UpCount),
			Down:     int(v.DownCount),
			Datetime: v.DateTime,
		}
		newKLines = append(newKLines, kline)
	}
	if len(newKLines) > 0 {
		// 复权之前, 假定当前缓存之中的数据都是复权过的数据
		// 那么就应该只拉取缓存最后1条记录之后的除权除息记录进行复权
		// 前复权adjust
		xdxrs := GetCacheXdxrList(securityCode)
		cacheLastDay := newKLines[len(newKLines)-1].Date
		for i := 0; i < len(xdxrs); i++ {
			xdxr := xdxrs[i]
			if xdxr.Category != 1 || xdxr.Date < startDate || xdxr.Date > cacheLastDay {
				// 忽略非除权信息以及除权数据在新数据之前的除权记录
				continue
			}
			xdxrDate := xdxr.Date
			factor := xdxr.Adjust()
			for j := 0; j < len(newKLines); j++ {
				kl := &newKLines[j]
				barCurrentDate := kl.Date
				if barCurrentDate > xdxrDate {
					break
				}
				if barCurrentDate < xdxrDate {
					kl.Open = factor(kl.Open)
					kl.Close = factor(kl.Close)
					kl.High = factor(kl.High)
					kl.Low = factor(kl.Low)
					// 成交量复权
					// 1. 计算均价线
					maPrice := kl.Amount / kl.Volume
					// 2. 均价线复权
					maPrice = factor(maPrice)
					// 3. 以成交金额为基准, 用复权均价线计算成交量
					kl.Volume = kl.Amount / maPrice
				}
				//plc := m1["last_close"].(reflect.Value)
				//plc.SetFloat(fuquan(plc.Float()))
				if barCurrentDate == xdxrDate {
					break
				}
			}
		}
	}
	var klines []KLine
	if kLength > klineDaysOffset {
		klines = cacheKLines[:kLength-klineDaysOffset]
	}
	if len(klines) > 0 {
		klines = append(klines, newKLines...)
	} else {
		klines = newKLines
	}
	if len(klines) > 0 {
		UpdateCacheKLines(securityCode, klines)
		fname := cache.KLineFilename(securityCode)
		_ = api.SlicesToCsv(fname, klines)
	}
	return klines
}
