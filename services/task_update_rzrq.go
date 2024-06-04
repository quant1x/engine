package services

import (
	"gitee.com/quant1x/engine/cache"
	"gitee.com/quant1x/engine/factors"
	"gitee.com/quant1x/engine/global"
	"gitee.com/quant1x/engine/market"
	"gitee.com/quant1x/gox/logger"
)

func jobUpdateMarginTrading() {
	logger.Infof("同步融资融券...")
	variables := global.GetGlobalVariables()
	updateMarginTrading(*variables.MarketData)
	logger.Infof("同步融资融券...OK")
}

func updateMarginTrading(marketData market.MarketData) {
	date := cache.DefaultCanReadDate()
	factors.MarginTradingTargetInit(date)
	allCodes := marketData.GetCodeList()
	for _, securityCode := range allCodes {
		misc := factors.GetL5Misc(securityCode)
		if misc == nil {
			continue
		}
		rzrq, ok := factors.GetMarginTradingTarget(securityCode)
		if ok {
			misc.RZYEZB = rzrq.RZYEZB
			factors.UpdateL5Misc(misc)
		}
	}
	factors.RefreshL5Misc()
}
