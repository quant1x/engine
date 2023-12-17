package services

import (
	"gitee.com/quant1x/engine/cache"
	"gitee.com/quant1x/engine/config"
	"gitee.com/quant1x/engine/market"
	"gitee.com/quant1x/engine/models"
	"gitee.com/quant1x/engine/realtime"
	"gitee.com/quant1x/engine/smart"
	"gitee.com/quant1x/engine/storages"
	"gitee.com/quant1x/engine/trader"
	"gitee.com/quant1x/gotdx/proto"
	"gitee.com/quant1x/gotdx/trading"
	"gitee.com/quant1x/gox/api"
	"gitee.com/quant1x/gox/logger"
	"gitee.com/quant1x/gox/num"
	"gitee.com/quant1x/gox/runtime"
	"slices"
)

// 任务 - 卖出117
func jobOneSizeFitsAllSales() {
	funcName, _, _ := runtime.Caller()
	updateInRealTime, status := trading.CanUpdateInRealtime()
	if updateInRealTime && IsTrading(status) {
		cookieCutterSell()
	} else {
		if runtime.Debug() {
			cookieCutterSell()
		} else {
			logger.Infof("%s, 非交易时段: %d", funcName, status)
		}
	}
}

// 一刀切卖出
func cookieCutterSell() {
	sellStrategyCode := models.ModelOneSizeFitsAllSells
	// 1. 获取117号策略(卖出)
	sellRule := config.GetTradeRule(sellStrategyCode)
	if sellRule == nil {
		return
	}
	// 2. 判断是否可以指定自动卖出
	if !sellRule.IsCookieCutterForSell() {
		return
	}
	// 3. 判断是否交易日
	if !trading.DateIsTradingDay() {
		return
	}
	// 3.1 判断是否交易时段
	if !sellRule.Session.IsTrading() {
		return
	}
	// 4. 查询持仓可卖的股票
	positions, err := trader.QueryHolding()
	if err != nil {
		return
	}
	// 5. 确定最早的持仓日期
	//firstDate := getEarlierDate(sellRule.HoldingPeriod)
	finalCodeList := checkoutCanSellStockList(sellStrategyCode)
	// 6. 遍历持仓
	direction := trader.SELL
	strategyName := sellRule.QmtStrategyName()
	for _, position := range positions {
		orderRemark := sellRule.Flag
		isNeedToSell := false
		// 6.1 如果持仓可卖数据小于1, 继续下一条持仓记录
		if position.CanUseVolume < 1 {
			continue
		}
		// 6.2 对齐证券代码
		stockCode := position.StockCode
		securityCode := proto.CorrectSecurityCode(stockCode)
		// 6.3 获取快照
		snapshot := models.GetTickFromMemory(securityCode)
		if snapshot == nil {
			continue
		}
		// 6.4 现价
		lastPrice := num.Decimal(snapshot.Price)
		// 昨日收盘
		lastClose := num.Decimal(snapshot.LastClose)
		// 6.5 计算涨停价
		limitUp, _ := market.PriceLimit(securityCode, lastClose)
		// 6.6 如果涨停, 则不出
		if lastPrice >= limitUp {
			continue
		}
		// 6.7 持仓成本
		avgPrice := position.OpenPrice
		// 6.8 盈亏比
		floatProfitLossRatio := num.NetChangeRate(avgPrice, lastPrice)
		_ = floatProfitLossRatio
		// 6.9 确定是否规则内最后一天持股
		isFinal := slices.Contains(finalCodeList, securityCode)
		// 117. 最后一天持股, 且是最后一个交易时段, 则卖出
		if isFinal && sellRule.Session.IsTodayLastSession() {
			// 卖出
			isNeedToSell = true
			orderRemark = "LASTDAY:"
			if floatProfitLossRatio > 0 {
				orderRemark = orderRemark + ">0"
			} else if floatProfitLossRatio == 0 {
				orderRemark = orderRemark + "=0"
			} else {
				orderRemark = orderRemark + "<0"
			}
		}
		// 6.10 股价高于前面一天最高价，且大于等于5日线，如果是获利的，卖出
		if !isNeedToSell {
			// 6.10.1 获取历史特征数据
			history := smart.GetL5History(securityCode)
			if history == nil {
				continue
			}
			// 6.10.2 计算5日均线
			ma5 := realtime.IncrementalMovingAverage(history.MA4, 5, lastPrice)
			// 6.10.3 股价高于前一天最高价, 且站上5日线以及盈利的情况下
			if lastPrice > history.HIGH && lastPrice >= ma5 && floatProfitLossRatio > 0 {
				// 卖出
				isNeedToSell = true
				orderRemark = "PL:P>H>MA5>0"
			}
		}
		// 18168
		if !isNeedToSell {
			continue
		}
		// 成本价
		//costPrice := position.OpenPrice
		orderPrice := lastPrice
		orderVolume := position.CanUseVolume
		// 卖出
		order_id, err := trader.DirectOrder(direction, strategyName, orderRemark, securityCode, orderPrice, orderVolume)
		if err != nil {
			continue
		}
		_ = order_id
	}
}

// 获得T+HoldingPeriod的具体日期
func getEarlierDate(period int) string {
	dates := trading.LastNDate(trading.Today(), period)
	earlier_date := trading.FixTradeDate(dates[0], cache.CACHE_DATE)
	return earlier_date
}

// 获取所有挂接了指定的卖出策略ID的交易规则
func getTradeRuleList(sellStrategyId int) []config.TradeRule {
	traderConfig := config.TraderConfig()
	var list []config.TradeRule
	for _, v := range traderConfig.Strategies {
		if v.Flag == models.OrderFlagSell || v.SellStrategy != sellStrategyId {
			continue
		}
		list = append(list, v)
	}
	return list
}

// 捡出T+HoldingPeriod日的股票列表
func checkoutCanSellStockList(sellStrategyId int) []string {
	var list []string
	tradeRules := getTradeRuleList(sellStrategyId)
	if len(tradeRules) == 0 {
		return list
	}
	for _, v := range tradeRules {
		date := getEarlierDate(v.HoldingPeriod)
		qmtStrategyName := v.QmtStrategyName()
		codes := storages.FetchListForFirstPurchase(date, qmtStrategyName, trader.BUY)
		if len(codes) == 0 {
			continue
		}
		list = append(list, codes...)
	}
	list = api.Unique(list)
	return list
}
