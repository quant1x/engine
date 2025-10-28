package services

import (
	"slices"
	"strings"

	"gitee.com/quant1x/data/exchange"
	"gitee.com/quant1x/engine/cache"
	"gitee.com/quant1x/engine/config"
	"gitee.com/quant1x/engine/models"
	"gitee.com/quant1x/engine/storages"
	"gitee.com/quant1x/engine/trader"
	"gitee.com/quant1x/gox/api"
	"gitee.com/quant1x/gox/logger"
)

// 获得T+HoldingPeriod的具体日期
func getEarlierDate(period int) string {
	dates := exchange.LastNDate(exchange.LastTradeDate(), period)
	earlier_date := exchange.FixTradeDate(dates[0], cache.CACHE_DATE)
	return earlier_date
}

// 不包含最后一个交易日的持股日期列表
func getHoldingDates(period int) []string {
	dates := exchange.LastNDate(exchange.LastTradeDate(), period)
	for i := 0; i < len(dates); i++ {
		dates[i] = exchange.FixTradeDate(dates[i], cache.CACHE_DATE)
	}
	return dates
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
func CheckoutCanSellStockList(sellStrategyId uint64, holdings []string) []string {
	tradeRules := getStrategyParameterList(sellStrategyId)
	if len(tradeRules) == 0 {
		return nil
	}
	// 1. 到期
	var listCanSell []string
	for _, v := range tradeRules {
		dates := getHoldingDates(v.HoldingPeriod)
		// 1. 到期日
		earlierDate := dates[0]
		qmtStrategyName := v.QmtStrategyName()
		codes := storages.FetchListForFirstPurchase(earlierDate, qmtStrategyName, trader.BUY)
		logger.Infof("sell strategy[%d]: from %d, last-day codes=%s", sellStrategyId, v.Id, strings.Join(codes, ","))
		if len(codes) == 0 {
			continue
		}
		// 筛选包含持股到期日的个股
		codes = api.Filter(codes, func(s string) bool {
			return slices.Contains(holdings, s)
		})
		listCanSell = append(listCanSell, codes...)
	}
	listCanSell = api.Unique(listCanSell)
	// 2. 未到期
	var listNoSell []string
	for _, v := range tradeRules {
		dates := getHoldingDates(v.HoldingPeriod)
		qmtStrategyName := v.QmtStrategyName()
		// 2. 未到期
		for _, orderDate := range dates[1:] {
			codes := storages.FetchListForFirstPurchase(orderDate, qmtStrategyName, trader.BUY)
			// 剔除包含持股到期日的个股
			codes = api.Filter(codes, func(s string) bool {
				return !slices.Contains(listCanSell, s)
			})
			listNoSell = append(listNoSell, codes...)
		}
	}
	listNoSell = api.Unique(listNoSell)
	// 3. 矫正持仓过期未卖出的个股
	for _, code := range holdings {
		if slices.Contains(listCanSell, code) || slices.Contains(listNoSell, code) {
			// 到期和未到期的忽略
			continue
		}
		// 流程走这里, 一般是非机器自动交易或者前一个交易日自动没出未成交造成的
		listCanSell = append(listCanSell, code)
	}
	return listCanSell
}

//// CheckoutUnsellableStockList 捡出不可卖的股票列表, T+HoldingPeriod日内的股票列表
//func CheckoutUnsellableStockList(sellStrategyId uint64) []string {
//	return nil
//}
