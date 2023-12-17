package models

import (
	"gitee.com/quant1x/engine/smart"
	"gitee.com/quant1x/gotdx"
	"gitee.com/quant1x/gotdx/quotes"
	"gitee.com/quant1x/gotdx/securities"
	"gitee.com/quant1x/gotdx/trading"
	"gitee.com/quant1x/gox/api"
	"gitee.com/quant1x/gox/logger"
	"gitee.com/quant1x/gox/num"
	"gitee.com/quant1x/gox/progressbar"
	"sync"
)

var (
	__mutexTicks sync.RWMutex
	__cacheTicks = map[string]quotes.Snapshot{}
)

// GetTickFromMemory 获取快照缓存
func GetTickFromMemory(securityCode string) *quotes.Snapshot {
	__mutexTicks.RLock()
	defer __mutexTicks.RUnlock()
	v, found := __cacheTicks[securityCode]
	if found {
		return &v
	}
	return nil
}

// GetStrategySnapshot 从缓存中获取快照
func GetStrategySnapshot(securityCode string) *QuoteSnapshot {
	v := GetTickFromMemory(securityCode)
	if v == nil || v.State != quotes.TDX_SECURITY_TRADE_STATE_NORMAL {
		// 非正常交易的记录忽略掉
		return nil
	}
	snapshot := QuoteSnapshot{}
	_ = api.Copy(&snapshot, &v)
	snapshot.Name = securities.GetStockName(securityCode)
	snapshot.Code = securityCode
	snapshot.OpeningChangeRate = num.NetChangeRate(snapshot.LastClose, snapshot.Open)
	snapshot.ChangeRate = num.NetChangeRate(snapshot.LastClose, snapshot.Price)
	f10 := smart.GetL5F10(securityCode)
	if f10 != nil {
		snapshot.Capital = f10.Capital
		snapshot.FreeCapital = f10.FreeCapital
		snapshot.OpenTurnZ = 10000 * float64(snapshot.OpenVolume) / float64(snapshot.FreeCapital)
	}
	history := smart.GetL5History(securityCode)
	if history != nil {
		lastMinuteVolume := history.GetMV5()
		snapshot.OpenQuantityRatio = float64(snapshot.OpenVolume) / lastMinuteVolume
		minuteVolume := float64(snapshot.Vol) / float64(trading.Minutes(snapshot.Date))
		snapshot.QuantityRatio = minuteVolume / lastMinuteVolume
	}
	snapshot.OpenBiddingDirection, snapshot.OpenVolumeDirection = v.CheckDirection()
	return &snapshot
}

// SyncAllSnapshots 同步快照数据
func SyncAllSnapshots(barIndex *int) {
	modName := "同步快照数据"
	allCodes := securities.AllCodeList()
	count := len(allCodes)
	var bar *progressbar.Bar = nil
	if barIndex != nil {
		bar = progressbar.NewBar(*barIndex, "执行["+modName+"]", count)
	}
	currentDate := trading.GetCurrentlyDay()
	tdxApi := gotdx.GetTdxApi()
	var snapshots []quotes.Snapshot
	for start := 0; start < count; start += quotes.TDX_SECURITY_QUOTES_MAX {
		length := count - start
		if length >= quotes.TDX_SECURITY_QUOTES_MAX {
			length = quotes.TDX_SECURITY_QUOTES_MAX
		}
		var subCodes []string
		for i := 0; i < length; i++ {
			securityCode := allCodes[start+i]
			subCodes = append(subCodes, securityCode)
			if barIndex != nil {
				bar.Add(1)
			}
		}
		if len(subCodes) == 0 {
			continue
		}
		for i := 0; i < quotes.DefaultRetryTimes; i++ {
			list, err := tdxApi.GetSnapshot(subCodes)
			if err != nil {
				logger.Errorf("ZS: 网络异常: %+v, 重试: %d", err, i+1)
				continue
			}
			for _, v := range list {
				// 修订日期
				v.Date = currentDate
				snapshots = append(snapshots, v)
			}
			break
		}
	}
	__mutexTicks.Lock()
	for _, v := range snapshots {
		__cacheTicks[v.SecurityCode] = v
	}
	__mutexTicks.Unlock()
	if barIndex != nil {
		*barIndex++
	}
}
