package services

import (
	"github.com/quant1x/engine/cache"
	"github.com/quant1x/engine/factors"
	"github.com/quant1x/engine/market"
	"github.com/quant1x/x/api"
	"github.com/quant1x/x/logger"
)

func jobUpdateMarginTrading() {
	logger.Infof("同步融资融券...")
	date := cache.DefaultCanReadDate()
	factors.MarginTradingTargetInit(date)
	updateMarginTradingForMisc(date)
	updateMarginTradingForRzrq(date)
	logger.Infof("同步融资融券...OK")
}

func updateMarginTradingForMisc(cacheDate string) {
	allCodes := market.GetCodeList()
	for _, securityCode := range allCodes {
		misc := factors.GetL5Misc(securityCode, cacheDate)
		if misc == nil {
			continue
		}
		rzrq, ok := factors.GetMarginTradingTarget(securityCode)
		if ok {
			misc.RZYEZB = rzrq.RZYEZB
			misc.UpdateTime = factors.GetTimestamp()
			factors.UpdateL5Misc(misc)
		}
	}
	factors.RefreshL5Misc()
}

func updateMarginTradingForRzrq(cacheDate string) {
	allCodes := market.GetCodeList()
	for _, securityCode := range allCodes {
		smt := factors.GetL5SecuritiesMarginTrading(securityCode, cacheDate)
		if smt == nil {
			continue
		}
		rzrq, ok := factors.GetMarginTradingTarget(securityCode)
		if ok {
			_ = api.Copy(smt, &rzrq)
			smt.UpdateTime = factors.GetTimestamp()
			smt.State |= factors.FeatureSecuritiesMarginTrading
			factors.UpdateL5SecuritiesMarginTrading(smt)
		}
	}
	factors.RefreshL5SecuritiesMarginTrading()
}
