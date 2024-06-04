package market

import (
	"gitee.com/quant1x/engine/config"
	"gitee.com/quant1x/engine/db"
	"gitee.com/quant1x/exchange"
	"gitee.com/quant1x/gotdx/securities"
	"gitee.com/quant1x/num"
	"strings"
)

type MarketData struct {
	db *db.Database
}

func NewMarketData(db *db.Database) *MarketData {
	return &MarketData{db: db}
}

// GetStockCodeList 获取股票代码列表。
func (m *MarketData) GetStockCodeList() []string {
	return m.db.GetStockCodeList(true)
}

// StockList 取得可以交易的证券代码列表
func (m *MarketData) GetStrategyStockCodeList(sp *config.StrategyParameter) []string {
	var codes []string
	for _, v := range sp.Sectors {
		sectorCode := strings.TrimSpace(v)
		if !strings.HasPrefix(sectorCode, config.GetSectorIgnorePrefix()) {
			blockInfo := securities.GetBlockInfo(sectorCode)
			if blockInfo != nil {
				codes = append(codes, blockInfo.ConstituentStocks...)

			}
		}
	}
	if len(codes) == 0 {
		codes = m.GetStockCodeList()
	}
	codes = sp.Filter(codes)
	return codes
}

// GetCodeList 加载所有股票代码，包括指数代码和板块代码。
func (m *MarketData) GetCodeList() []string {
	allCodes := make([]string, 0)

	// 添加指数代码。
	allCodes = append(allCodes, exchange.IndexList()...)

	// 添加板块代码。
	blocks := securities.BlockList()
	for _, v := range blocks {
		allCodes = append(allCodes, v.Code)
	}

	// 添加股票代码。
	stockCodes := m.GetStockCodeList()
	allCodes = append(allCodes, stockCodes...)

	return allCodes
}

// PriceLimit 计算基于前一个收盘价的涨停和跌停价格。
func (m *MarketData) PriceLimit(securityCode string, lastClose float64) (limitUp, limitDown float64) {
	limitRate := exchange.MarketLimit(securityCode)
	priceLimitUp := num.Decimal(lastClose * (1.000 + limitRate))
	priceLimitDown := num.Decimal(lastClose * (1.000 - limitRate))
	return priceLimitUp, priceLimitDown
}
