package models

import (
	"gitee.com/quant1x/gotdx"
	"gitee.com/quant1x/gotdx/proto"
	"gitee.com/quant1x/gotdx/quotes"
	"gitee.com/quant1x/gotdx/securities"
	"gitee.com/quant1x/gotdx/trading"
	"gitee.com/quant1x/gox/logger"
	"gitee.com/quant1x/gox/progressbar"
)

var (
	// 当日缓存数据
	__cacheSnapshots = map[string]quotes.Snapshot{}
)

func getQuoteSnapshot(securityCode string) *quotes.Snapshot {
	v, found := __cacheSnapshots[securityCode]
	if found {
		return &v
	}
	return nil
}

//// GetStrategySnapshot 从缓存中获取快照
//func GetStrategySnapshot(securityCode string) *models.QuoteSnapshot {
//	v := __cacheSnapshots[securityCode]
//	if v.State != quotes.TDX_SECURITY_TRADE_STATE_NORMAL {
//		// 非正常交易的记录忽略掉
//		return nil
//	}
//	snapshot := models.QuoteSnapshot{}
//	_ = api.Copy(&snapshot, &v)
//	snapshot.Name = securities.GetStockName(securityCode)
//	snapshot.Code = securityCode
//	snapshot.OpeningChangeRate = num.NetChangeRate(snapshot.LastClose, snapshot.Open)
//	snapshot.ChangeRate = num.NetChangeRate(snapshot.LastClose, snapshot.Price)
//	f10 := flash.GetL5F10(securityCode)
//	if f10 != nil {
//		snapshot.Capital = f10.Capital
//		snapshot.FreeCapital = f10.FreeCapital
//		snapshot.OpenTurnZ = 10000 * float64(snapshot.OpenVolume) / float64(snapshot.FreeCapital)
//	}
//	extend := flash.GetL5Exchange(securityCode)
//	if extend != nil {
//		snapshot.QuantityRatio = float64(snapshot.OpenVolume) / extend.MAV5
//	}
//	snapshot.OpenBiddingDirection, snapshot.OpenVolumeDirection = v.CheckDirection()
//	return &snapshot
//}

// GetAllSnapshots 同步快照数据
func GetAllSnapshots(barIndex *int) {
	modName := "同步快照数据"
	allCodes := securities.AllCodeList()
	count := len(allCodes)
	bar := progressbar.NewBar(*barIndex, "执行["+modName+"]", count)
	for start := 0; start < count; start += quotes.TDX_SECURITY_QUOTES_MAX {
		length := count - start
		if length >= quotes.TDX_SECURITY_QUOTES_MAX {
			length = quotes.TDX_SECURITY_QUOTES_MAX
		}
		subCodes := []string{}
		for i := 0; i < length; i++ {
			securityCode := allCodes[start+i]
			subCodes = append(subCodes, securityCode)
			bar.Add(1)
		}
		if len(subCodes) == 0 {
			continue
		}
		tdxApi := gotdx.GetTdxApi()
		currentDate := trading.GetCurrentlyDay()
		for i := 0; i < quotes.DefaultRetryTimes; i++ {
			list, err := tdxApi.GetSnapshot(subCodes)
			if err != nil {
				logger.Errorf("ZS: 网络异常: %+v, 重试: %d", err, i+1)
				continue
			}
			for _, v := range list {
				// 修订日期
				v.Date = currentDate
				securityCode := proto.GetSecurityCode(v.Market, v.Code)
				v.Code = securityCode
				__cacheSnapshots[v.Code] = v
			}
			break
		}
	}
	*barIndex++
}