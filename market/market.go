package market

import (
	"fmt"

	"gitee.com/quant1x/data/exchange"
	"gitee.com/quant1x/data/level1/securities"
	"gitee.com/quant1x/num"
)

func GetStockCodeList() []string {
	var allCodes []string
	// 上海
	// sh600000-sh609999
	{
		var (
			codeBegin = 600000
			codeEnd   = 609999
		)
		for i := codeBegin; i <= codeEnd; i++ {
			fc := fmt.Sprintf("sh%d", i)
			if IsNeedIgnore(fc) {
				continue
			}
			allCodes = append(allCodes, fc)
		}
	}
	// 科创板
	// sh688000-sh688999
	{
		var (
			codeBegin = 688000
			codeEnd   = 689999
		)
		for i := codeBegin; i <= codeEnd; i++ {
			fc := fmt.Sprintf("sh%d", i)
			if IsNeedIgnore(fc) {
				continue
			}
			allCodes = append(allCodes, fc)
		}
	}
	// 深圳证券交易所
	// 深圳主板: sz000000-sz000999
	{
		var (
			codeBegin = 0
			codeEnd   = 999
		)
		for i := codeBegin; i <= codeEnd; i++ {
			fc := fmt.Sprintf("sz000%03d", i)
			if IsNeedIgnore(fc) {
				continue
			}
			allCodes = append(allCodes, fc)
		}
	}
	// 中小板: sz001000-sz009999
	{
		var (
			codeBegin = 1000
			codeEnd   = 9999
		)
		for i := codeBegin; i <= codeEnd; i++ {
			fc := fmt.Sprintf("sz00%04d", i)
			if IsNeedIgnore(fc) {
				continue
			}
			allCodes = append(allCodes, fc)
		}
	}
	// 创业板: sz300000-sz300999
	{
		var (
			codeBegin = 300000
			codeEnd   = 309999
		)
		for i := codeBegin; i <= codeEnd; i++ {
			fc := fmt.Sprintf("sz%06d", i)
			if IsNeedIgnore(fc) {
				continue
			}
			allCodes = append(allCodes, fc)
		}
	}

	// 北交所: bj920000-bj920999
	{
		var (
			codeBegin = 920000
			codeEnd   = 920999
		)
		for i := codeBegin; i <= codeEnd; i++ {
			fc := fmt.Sprintf("bj%06d", i)
			if IsNeedIgnore(fc) {
				continue
			}
			allCodes = append(allCodes, fc)
		}
	}

	//allCodes = allCodes[0:0]
	// 港股: hk00001-hk09999
	//{
	//	var (
	//		codeBegin = 1
	//		codeEnd   = 9999
	//	)
	//	for i := codeBegin; i <= codeEnd; i++ {
	//		fc := fmt.Sprintf("hk%05d", i)
	//		allCodes = append(allCodes, fc)
	//	}
	//}

	return allCodes
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
