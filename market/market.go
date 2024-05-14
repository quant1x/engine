package market

import (
	"gitee.com/quant1x/engine/storages"
	"gitee.com/quant1x/exchange"
	"gitee.com/quant1x/gotdx/securities"
	"gitee.com/quant1x/num"
)

func GetStockCodeList() []string {
	return storages.GlobalDB.GetStockCodeList(true)
}

// GetCodeList 加载全部股票代码
func GetCodeList() []string {
	allCodes := make([]string, 0)
	// 添加指数代码
	allCodes = append(allCodes, exchange.IndexList()...)

	// 板块信息
	blocks := securities.BlockList()
	for _, v := range blocks {
		allCodes = append(allCodes, v.Code)
	}
	stockCodes := GetStockCodeList()
	allCodes = append(allCodes, stockCodes...)
	return allCodes
}

// PriceLimit 计算涨停板和跌停板的价格
func PriceLimit(securityCode string, lastClose float64) (limitUp, limitDown float64) {
	limitRate := exchange.MarketLimit(securityCode)
	priceLimitUp := num.Decimal(lastClose * (1.000 + limitRate))
	priceLimitDown := num.Decimal(lastClose * (1.000 - limitRate))
	return priceLimitUp, priceLimitDown
}
