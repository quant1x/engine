package services

import (
	"fmt"
	"gitee.com/quant1x/engine/config"
	"gitee.com/quant1x/engine/market"
	"gitee.com/quant1x/engine/models"
	"gitee.com/quant1x/engine/trader"
	"gitee.com/quant1x/gotdx/proto"
	"gitee.com/quant1x/gotdx/trading"
	"gitee.com/quant1x/gox/num"
)

// 一刀切卖出
func cookieCutterSell() {
	tradeRule := config.GetTradeRule(models.ModelOneSizeFitsAllSales)
	// 1. 判断是否可以指定自动卖出
	if !tradeRule.IsCookieCutterForSell() {
		return
	}
	// 2. 判断是否交易日
	if !trading.DateIsTradingDay() {
		return
	}
	// 2. 查询持仓可卖的股票
	positions, err := trader.QueryHolding()
	if err != nil {
		return
	}
	// 3. 遍历持仓
	direction := trader.SELL
	strategyName := fmt.Sprintf("S%d", tradeRule.Id)
	orderRemark := tradeRule.Flag
	for _, position := range positions {
		if position.CanUseVolume < 1 {
			continue
		}
		stockCode := position.StockCode
		securityCode := proto.CorrectSecurityCode(stockCode)
		// 获取快照
		snapshot := models.GetTickFromMemory(securityCode)
		if snapshot == nil {
			continue
		}
		// 现价
		lastPrice := num.Decimal(snapshot.Price)
		// 昨日收盘
		lastClose := num.Decimal(snapshot.LastClose)
		// 计算涨停价
		limitUp, _ := market.PriceLimit(securityCode, lastClose)
		// 如果涨停, 不出
		if lastPrice >= limitUp {
			continue
		}
		// 成本价
		//costPrice := position.OpenPrice
		orderPrice := lastPrice
		orderVolume := position.CanUseVolume
		order_id, err := trader.DirectOrder(direction, strategyName, orderRemark, securityCode, orderPrice, orderVolume)
		if err != nil {
			continue
		}
		_ = order_id
	}
}
