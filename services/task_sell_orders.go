package services

import (
	"gitee.com/quant1x/engine/cache"
	"gitee.com/quant1x/engine/config"
	"gitee.com/quant1x/engine/models"
	"gitee.com/quant1x/engine/storages"
	"gitee.com/quant1x/engine/trader"
	"gitee.com/quant1x/exchange"
	"gitee.com/quant1x/gox/api"
	"gitee.com/quant1x/gox/logger"
	"strings"
)

// 获得T+HoldingPeriod的具体日期
func getEarlierDate(period int) string {
	dates := exchange.LastNDate(exchange.LastTradeDate(), period)
	earlier_date := exchange.FixTradeDate(dates[0], cache.CACHE_DATE)
	return earlier_date
}

// 获取所有挂接了指定的卖出策略ID的交易规则
func getStrategyParameterList(sellStrategyId uint64) []config.StrategyParameter {
	traderConfig := config.TraderConfig()
	var list []config.StrategyParameter
	for _, v := range traderConfig.Strategies {
		if v.Flag == models.OrderFlagSell || v.SellStrategy != sellStrategyId {
			continue
		}
		list = append(list, v)
	}
	return list
}

// CheckoutCanSellStockList 捡出T+HoldingPeriod日的股票列表
func CheckoutCanSellStockList(sellStrategyId uint64) []string {
	var list []string
	tradeRules := getStrategyParameterList(sellStrategyId)
	if len(tradeRules) == 0 {
		return list
	}
	for _, v := range tradeRules {
		date := getEarlierDate(v.HoldingPeriod)
		qmtStrategyName := v.QmtStrategyName()
		codes := storages.FetchListForFirstPurchase(date, qmtStrategyName, trader.BUY)
		logger.Infof("sell strategy[%d]: from %d, last-day codes=%s", sellStrategyId, v.Id, strings.Join(codes, ","))
		if len(codes) == 0 {
			continue
		}
		list = append(list, codes...)
	}
	list = api.Unique(list)
	return list
}
