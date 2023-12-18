package models

import (
	"gitee.com/quant1x/engine/factors"
	"gitee.com/quant1x/engine/smart"
	"gitee.com/quant1x/gotdx"
	"gitee.com/quant1x/gotdx/proto"
	"gitee.com/quant1x/gotdx/quotes"
	"gitee.com/quant1x/gotdx/trading"
	"gitee.com/quant1x/gox/api"
	"gitee.com/quant1x/gox/logger"
	"gitee.com/quant1x/gox/num"
)

func QuoteSnapshotFromProtocol(v quotes.Snapshot) factors.QuoteSnapshot {
	snapshot := factors.QuoteSnapshot{}
	_ = api.Copy(&snapshot, &v)
	securityCode := proto.GetSecurityCode(v.Market, v.Code)
	snapshot.Code = securityCode
	snapshot.OpeningChangeRate = num.NetChangeRate(snapshot.LastClose, snapshot.Open)
	snapshot.ChangeRate = num.NetChangeRate(snapshot.LastClose, snapshot.Price)
	snapshot.PremiumRate = num.NetChangeRate(snapshot.Open, snapshot.Price)
	snapshot.OpenBiddingDirection, snapshot.OpenVolumeDirection = v.CheckDirection()
	// 涨跌力度
	snapshot.ChangePower = float64(snapshot.OpenVolume) / snapshot.OpeningChangeRate
	snapshot.AverageBiddingVolume = v.AverageBiddingVolume()

	// 补全F10相关
	f10 := smart.GetL5F10(securityCode)
	if f10 != nil {
		snapshot.Name = f10.SecurityName
		snapshot.Capital = f10.Capital
		snapshot.FreeCapital = f10.FreeCapital
		snapshot.OpenTurnZ = f10.TurnZ(snapshot.OpenVolume)
	}
	// 补全扩展相关
	history := smart.GetL5History(securityCode)
	if history != nil && history.MV5 > 0 {
		lastMinuteVolume := history.GetMV5()
		snapshot.OpenQuantityRatio = float64(snapshot.OpenVolume) / lastMinuteVolume
		minuteVolume := float64(snapshot.Vol) / float64(trading.Minutes(snapshot.Date))
		snapshot.QuantityRatio = minuteVolume / lastMinuteVolume
	}
	return snapshot
}

// BatchSnapShot 批量获取即时行情数据快照
func BatchSnapShot(codes []string) []factors.QuoteSnapshot {
	tdxApi := gotdx.GetTdxApi()
	list := []factors.QuoteSnapshot{}
	var err error
	var hq []quotes.Snapshot
	retryTimes := 0
	for retryTimes < quotes.DefaultRetryTimes {
		hq, err = tdxApi.GetSnapshot(codes)
		if err == nil && hq != nil {
			break
		}
		retryTimes++
	}

	if err != nil {
		logger.Errorf("获取即时行情数据失败", err)
		return list
	}

	for _, v := range hq {
		if v.State != quotes.TDX_SECURITY_TRADE_STATE_NORMAL {
			// 非正常交易的记录忽略掉
			continue
		}
		snapshot := QuoteSnapshotFromProtocol(v)
		list = append(list, snapshot)
	}
	return list
}

// GetTick 获取一只股票的tick数据
//
//	该函数用于测试
func GetTick(securityCode string) *factors.QuoteSnapshot {
	list := BatchSnapShot([]string{securityCode})
	if len(list) > 0 {
		v := list[0]
		return &v
	}
	return nil
}
