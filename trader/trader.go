package trader

import (
	"encoding/json"
	"gitee.com/quant1x/gox/http"
	"gitee.com/quant1x/gox/logger"
)

const (
	// miniQMT代理服务器地址
	urlPrefixMiniQmtProxy = "http://10.211.55.3:18168/qmt"
	// 查询前缀
	urlPrefixForQuery = urlPrefixMiniQmtProxy + "/query"
	// 交易前缀
	urlPrefixForTrade = urlPrefixMiniQmtProxy + "/trade"
	// 查询账户信息
	urlAccount = urlPrefixForQuery + "/asset"
	// 查询持仓信息
	urlHolding = urlPrefixForQuery + "/holding"
	// 查询委托
	urlOrders = urlPrefixForQuery + "/order"
)

// AccountDetail 账户信息
type AccountDetail struct {
	TotalAsset  float64 `name:"总金额" json:"total_asset"`
	Cash        float64 `name:"可用" json:"cash"`
	MarketValue float64 `name:"市值" json:"market_value"`
	FrozenCash  float64 `name:"冻结" json:"frozen_cash"`
}

// Position 持仓信息
type Position struct {
	StockCode    string  `name:"证券代码" json:"stock_code"`
	Volume       int     `name:"持仓量" json:"volume"`
	CanUseVolume int     `name:"可卖" json:"can_use_volume"`
	OpenPrice    float64 `name:"成本价" json:"open_price"`
	MarketValue  float64 `name:"市值" json:"market_value"`
}

// QueryAccount 查询账户信息
func QueryAccount() (*AccountDetail, error) {
	data, err := http.Post(urlAccount, "")
	if err != nil {
		logger.Errorf("trader: 查询账户异常: %+v", err)
		return nil, err
	}
	var detail AccountDetail
	err = json.Unmarshal(data, &detail)
	if err != nil {
		logger.Errorf("trader: 解析json异常: %+v", err)
		return nil, err
	}
	return &detail, nil
}

// QueryHolding 查询持仓
func QueryHolding() ([]Position, error) {
	data, err := http.Post(urlHolding, "")
	if err != nil {
		logger.Errorf("trader: 查询持仓异常: %+v", err)
		return nil, err
	}
	var detail []Position
	err = json.Unmarshal(data, &detail)
	if err != nil {
		logger.Errorf("trader: 解析json异常: %+v", err)
		return nil, err
	}
	return detail, nil
}

// QueryOrders 查询当日委托
func QueryOrders() {

}
