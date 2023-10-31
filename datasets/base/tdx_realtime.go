package base

import (
	"encoding/json"
	"fmt"
	"gitee.com/quant1x/engine/cache"
	"gitee.com/quant1x/gotdx"
	"gitee.com/quant1x/gotdx/proto"
	"gitee.com/quant1x/gotdx/quotes"
	"gitee.com/quant1x/gotdx/trading"
	"gitee.com/quant1x/gox/api"
	"gitee.com/quant1x/gox/errors"
	"gitee.com/quant1x/gox/logger"
	"os"
	"time"
)

// BatchRealtimeBasicKLine 批量获取实时行情数据
func BatchRealtimeBasicKLine(codes []string) error {
	if len(codes) > quotes.TDX_SECURITY_QUOTES_MAX {
		message := fmt.Sprintf("BatchKLineWideRealtime: codes count > %d", quotes.TDX_SECURITY_QUOTES_MAX)
		return errors.New(message)
	}
	now := time.Now()
	nowServerTime := now.Format(trading.CN_SERVERTIME_FORMAT)
	lastTradeday := trading.LastTradeDate()
	today := trading.Today()
	if lastTradeday != today {
		// 当天非交易日, 不更新, 直接返回
		return nil
	}
	if nowServerTime < trading.CN_TradingStartTime || nowServerTime > trading.CN_TradingStopTime {
		// 非交易时间, 不更新, 直接返回
		return nil
	}

	tdxApi := gotdx.GetTdxApi()
	var err error
	var hq []quotes.Snapshot
	retryTimes := 0
	for retryTimes < quotes.DefaultRetryTimes {
		list, err := tdxApi.GetSnapshot(codes)
		if err == nil && list != nil && len(list) > 0 {
			hq = list
			break
		}
		retryTimes++
	}
	if err != nil {
		logger.Errorf("获取即时行情数据失败", err)
		return err
	}
	for _, v := range hq {
		if v.State == quotes.TDX_SECURITY_TRADE_STATE_DELISTING || v.Code == proto.StockDelisting || v.LastClose == float64(0) {
			// 终止上市的数据略过
			continue
		}
		securityCode := proto.GetMarketFlag(v.Market) + v.Code
		kl := KLine{
			Date:   lastTradeday,
			Open:   v.Open,
			Close:  v.Price,
			High:   v.High,
			Low:    v.Low,
			Volume: float64(v.Vol),
			Amount: v.Amount,
			Up:     v.BVol,
			Down:   v.SVol,
		}
		cacheKLines := LoadBasicKline(securityCode)
		cacheLength := len(cacheKLines)
		if cacheLength == 0 {
			continue
		}
		// 获取缓存中最后一根K线的日期
		klineFilename := cache.KLineFilename(securityCode)
		cacheLastDate := cacheKLines[cacheLength-1].Date
		if len(cacheLastDate) == 0 {
			// 如果缓存文件异常, 则删除
			data, _ := json.Marshal(cacheKLines[cacheLength-1])
			text := api.Bytes2String(data)
			logger.Errorf("realtime kline error, code: %s, date=[%s]", securityCode, text)
			_ = os.Remove(klineFilename)
			// 全量更新K线
			UpdateAllBasicKLine(securityCode)
			continue
		}
		ts := trading.TradeRange(cacheLastDate, lastTradeday)
		if len(ts) > 2 {
			// 超过2天的差距, 不能用realtime更新K线数据
			// 只能是当天更新 或者是新增, 跨越2个以上的交易日不更新
			// 全量更新K线
			UpdateAllBasicKLine(securityCode)
			continue
		}
		// 数据差异数
		diffDays := 0
		// 当日的K线数据已经存在
		if cacheLastDate == lastTradeday {
			// 如果最后一条数据和最后一个交易日相同, 那么去掉缓存中的最后一条, 用实时数据填补
			// 这种情况的出现是K线被更新过了, 现在做的是用快照更新K线
			diffDays = 1
		} else if nowServerTime > v.ServerTime {
			diffDays = 0
		}
		var klines []KLine
		if diffDays > 0 {
			klines = cacheKLines[:cacheLength-diffDays]
		} else {
			klines = cacheKLines
		}
		// 连接缓存和实时数据
		klines = append(klines, kl)
		err := api.SlicesToCsv(klineFilename, klines)
		if err != nil {
			logger.Errorf("更新K线数据文件失败:%s", v.Code)
		}
	}
	return nil
}

// BasicKLineForSnapshot 通过snapshot更新基础K线
func BasicKLineForSnapshot(v quotes.Snapshot) {
	now := time.Now()
	nowServerTime := now.Format(trading.CN_SERVERTIME_FORMAT)
	lastTradeday := trading.LastTradeDate()
	today := trading.Today()
	if lastTradeday != today {
		// 当天非交易日, 不更新, 直接返回
		if !cache.Debug {
			return
		}
	}
	if nowServerTime < trading.CN_TradingStartTime || nowServerTime > trading.CN_TradingStopTime {
		// 非交易时间, 不更新, 直接返回
		if !cache.Debug {
			return
		}
	}
	if v.State == quotes.TDX_SECURITY_TRADE_STATE_DELISTING || v.Code == proto.StockDelisting || v.LastClose == float64(0) {
		// 终止上市的数据略过
		return
	}
	//securityCode := proto.GetMarketFlag(v.Market) + v.Code
	securityCode := v.SecurityCode
	kl := KLine{
		Date:   lastTradeday,
		Open:   v.Open,
		Close:  v.Price,
		High:   v.High,
		Low:    v.Low,
		Volume: float64(v.Vol),
		Amount: v.Amount,
		Up:     v.BVol,
		Down:   v.SVol,
	}
	cacheKLines := LoadBasicKline(securityCode)
	cacheLength := len(cacheKLines)
	if cacheLength == 0 {
		return
	}
	// 获取缓存中最后一根K线的日期
	klineFilename := cache.KLineFilename(securityCode)
	cacheLastDate := cacheKLines[cacheLength-1].Date
	if len(cacheLastDate) == 0 {
		// 如果缓存文件异常, 则删除
		data, _ := json.Marshal(cacheKLines[cacheLength-1])
		text := api.Bytes2String(data)
		logger.Errorf("realtime kline error, code: %s, date=[%s]", securityCode, text)
		_ = os.Remove(klineFilename)
		// 全量更新K线
		UpdateAllBasicKLine(securityCode)
		return
	}
	ts := trading.TradeRange(cacheLastDate, lastTradeday)
	if len(ts) > 2 {
		// 超过2天的差距, 不能用realtime更新K线数据
		// 只能是当天更新 或者是新增, 跨越2个以上的交易日不更新
		// 全量更新K线
		UpdateAllBasicKLine(securityCode)
		return
	}
	// 数据差异数
	diffDays := 0
	// 当日的K线数据已经存在
	if cacheLastDate == lastTradeday {
		// 如果最后一条数据和最后一个交易日相同, 那么去掉缓存中的最后一条, 用实时数据填补
		// 这种情况的出现是K线被更新过了, 现在做的是用快照更新K线
		diffDays = 1
	} else if nowServerTime > v.ServerTime {
		diffDays = 0
	}
	var klines []KLine
	if diffDays > 0 {
		klines = cacheKLines[:cacheLength-diffDays]
	} else {
		klines = cacheKLines
	}
	// 更新缓存K线
	UpdateCacheKLines(securityCode, klines)
	// 连接缓存和实时数据
	klines = append(klines, kl)
	err := api.SlicesToCsv(klineFilename, klines)
	if err != nil {
		logger.Errorf("更新K线数据文件失败:%s", v.Code)
	}
}
