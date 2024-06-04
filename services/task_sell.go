package services

import (
	"fmt"
	"gitee.com/quant1x/engine/config"
	"gitee.com/quant1x/engine/factors"
	"gitee.com/quant1x/engine/global"
	"gitee.com/quant1x/engine/market"
	"gitee.com/quant1x/engine/models"
	"gitee.com/quant1x/engine/realtime"
	"gitee.com/quant1x/engine/trader"
	"gitee.com/quant1x/exchange"
	"gitee.com/quant1x/gox/logger"
	"gitee.com/quant1x/gox/runtime"
	"gitee.com/quant1x/num"
	"slices"
)

// 任务 - 卖出117
func jobOneSizeFitsAllSales() {
	updateInRealTime, status := exchange.CanUpdateInRealtime()
	variables := global.GetGlobalVariables()
	if updateInRealTime && IsTrading(status) {
		cookieCutterSell(*variables.MarketData)
	} else if runtime.Debug() {
		cookieCutterSell(*variables.MarketData)
	}
}

// 一刀切卖出
func cookieCutterSell(marketData market.MarketData) {
	defer runtime.IgnorePanic("")
	sellStrategyCode := models.ModelOneSizeFitsAllSells
	// 1. 获取117号策略(卖出)
	sellRule := config.GetStrategyParameterByCode(sellStrategyCode)
	if sellRule == nil {
		return
	}
	// 2. 判断是否可以指定自动卖出
	if !sellRule.IsCookieCutterForSell() {
		return
	}
	// 3. 判断是否交易日
	if !exchange.DateIsTradingDay() {
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
	// 5. 确定持股到期的个股列表
	var holdings []string
	for _, position := range positions {
		if position.CanUseVolume < 1 {
			continue
		}
		stockCode := position.StockCode
		securityCode := exchange.CorrectSecurityCode(stockCode)
		holdings = append(holdings, securityCode)
	}
	finalCodeList := CheckoutCanSellStockList(sellStrategyCode, holdings)
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
		// 6.2 对齐证券代码, 检查黑白名单
		stockCode := position.StockCode
		securityCode := exchange.CorrectSecurityCode(stockCode)
		if !trader.CheckForSell(securityCode) {
			// 禁止卖出, 则返回
			logger.Infof("%s[%d]: %s ProhibitForBuying", sellRule.Name, sellRule.Id, securityCode)
			continue
		}
		// 6.3 获取快照
		snapshot := models.GetStrategySnapshot(securityCode)
		if snapshot == nil {
			continue
		}
		// 6.4 现价
		lastPrice := num.Decimal(snapshot.Price)
		// 昨日收盘
		lastClose := num.Decimal(snapshot.LastClose)
		// 6.5 计算涨停价
		limitUp, _ := marketData.PriceLimit(securityCode, lastClose)
		// 6.6 如果涨停, 则不出
		if lastPrice >= limitUp {
			logger.Infof("%s[%d]: %s LimitUp, skip", sellRule.Name, sellRule.Id, securityCode)
			continue
		}
		// 6.7 持仓成本
		avgPrice := position.OpenPrice
		// 6.8 盈亏比
		floatProfitLossRatio := num.NetChangeRate(avgPrice, lastPrice)
		// 6.9 确定是否规则内最后一天持股
		isFinal := slices.Contains(finalCodeList, securityCode)
		todayLastSession := sellRule.Session.IsTodayLastSession()
		logger.Infof("%s[%d]: %s, profit-loss-ratio=%.02f, last-day=%t, last-session=%t", sellRule.Name, sellRule.Id, securityCode, floatProfitLossRatio, isFinal, todayLastSession)
		// 117. 最后一天持股, 且是最后一个交易时段, 则卖出
		if isFinal && todayLastSession {
			// 卖出
			isNeedToSell = true
			orderRemark = "LASTDAY:P"
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
			history := factors.GetL5History(securityCode)
			if history == nil {
				continue
			}
			// 6.10.2 计算5日均线
			ma5 := realtime.IncrementalMovingAverage(history.MA4, 5, lastPrice)
			// 风险收益比（Risk/Reward Ratio）
			orderRemark = "RISK:P"
			// 6.10.3 股价高于前一天最高价, 且站上5日线以及盈利的情况下
			if lastPrice > history.HIGH && lastPrice >= ma5 && floatProfitLossRatio > 0 {
				// 卖出
				isNeedToSell = true
				orderRemark += ">H>MA5>0"
			} else {
				//6.11 如果股价触及止盈比例, 则卖出
				if sellRule.Session.CanTakeProfit() && floatProfitLossRatio > sellRule.TakeProfitRatio {
					isNeedToSell = true
					// 止盈
					orderRemark += ">TPR"
				} else if sellRule.Session.CanStopLoss() && floatProfitLossRatio < sellRule.StopLossRatio {
					isNeedToSell = true
					// 止损
					orderRemark += "<SLR"
				}
			}
		}
		// 7 卖出操作, 最后修订
		// 成本价
		//costPrice := position.OpenPrice
		orderPrice := lastPrice
		orderVolume := position.CanUseVolume
		// 7.1 如果跳空低开, 现价卖出, 如果想开盘就挂卖单, 就需要配置一个早盘集合竞价结束后, 早盘交易之前的交易时段
		// 比如 09:29:00~09:29:59
		if !isNeedToSell && snapshot.ExistDownwardGap() {
			isNeedToSell = true
			orderRemark = "GAP:DOWNWARD"
			orderPrice = lastPrice
		}
		// 7.2 如果卖出策略配置的固定收益率大于0, 则卖出
		if !isNeedToSell && sellRule.FixedYield > 0 {
			fee := trader.EvaluatePriceForSell(securityCode, avgPrice, orderVolume, sellRule.FixedYield)
			if fee != nil && fee.Price > avgPrice {
				isNeedToSell = true
				orderRemark = fmt.Sprintf("FIXEDYIELD:%.2f", sellRule.FixedYield*100)
				orderPrice = fee.Price
			}
		}
		// 18168
		if !isNeedToSell {
			continue
		}
		// 卖出
		order_id, err := trader.DirectOrder(direction, strategyName, orderRemark, securityCode, trader.LATEST_PRICE, orderPrice, orderVolume)
		if err != nil {
			continue
		}
		_ = order_id
	}
}
