package trader

import (
	"encoding/json"
	"fmt"
	"gitee.com/quant1x/engine/config"
	"gitee.com/quant1x/engine/models"
	"gitee.com/quant1x/gotdx/proto"
	"gitee.com/quant1x/gox/http"
	"gitee.com/quant1x/gox/logger"
	urlpkg "net/url"
	"strings"
)

var (
	// miniQMT代理服务器地址
	//urlPrefixMiniQmtProxy = "http://10.211.55.3:18168/qmt"
	urlPrefixMiniQmtProxy = config.TraderConfig().ProxyUrl
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
	// 委托
	urlPlaceOrder = urlPrefixForTrade + "/order"
	// 撤单
	urlCancelOrder = urlPrefixForTrade + "/cancel"
)

// Direction 交易方向
type Direction string

func (d Direction) String() string {
	return string(d)
}

const (
	BUY  Direction = "buy"  // 买入
	SELL Direction = "sell" // 卖出
)

type ProxyResult struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

// OrderResult 结果
type OrderResult struct {
	ProxyResult
	OrderId int `json:"order_id"`
}

// AccountDetail 账户信息
type AccountDetail struct {
	TotalAsset  float64 `name:"总金额" json:"total_asset"`
	Cash        float64 `name:"可用" json:"cash"`
	MarketValue float64 `name:"市值" json:"market_value"`
	FrozenCash  float64 `name:"冻结" json:"frozen_cash"`
}

// PositionDetail 持仓信息
type PositionDetail struct {
	StockCode    string  `name:"证券代码" json:"stock_code"`
	Volume       int     `name:"持仓量" json:"volume"`
	CanUseVolume int     `name:"可卖" json:"can_use_volume"`
	OpenPrice    float64 `name:"成本价" json:"open_price"`
	MarketValue  float64 `name:"市值" json:"market_value"`
}

// OrderDetail 委托订单
type OrderDetail struct {
	StockCode    string  `name:"证券代码" json:"stock_code"`
	OrderVolume  int     `name:"委托量" json:"order_volume"`
	TradedVolume int     `name:"成交量" json:"traded_volume"`
	Price        float64 `name:"委托价格" json:"price"`
	OrderType    int     `name:"订单类型" json:"order_type"`
	OrderStatus  int     `name:"订单状态" json:"order_status"`
	OrderId      int     `name:"订单ID" json:"order_id"`
	OrderSysid   string  `name:"合同编号" json:"order_sysid"`
	OrderTime    string  `name:"委托时间" json:"order_time"`
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
func QueryHolding() ([]PositionDetail, error) {
	data, err := http.Post(urlHolding, "")
	if err != nil {
		logger.Errorf("trader: 查询持仓异常: %+v", err)
		return nil, err
	}
	var detail []PositionDetail
	err = json.Unmarshal(data, &detail)
	if err != nil {
		logger.Errorf("trader: 解析json异常: %+v", err)
		return nil, err
	}
	return detail, nil
}

// QueryOrders 查询当日委托
func QueryOrders() ([]OrderDetail, error) {
	data, err := http.Post(urlOrders, "")
	if err != nil {
		logger.Errorf("trader: 查询委托异常: %+v", err)
		return nil, err
	}
	var detail []OrderDetail
	err = json.Unmarshal(data, &detail)
	if err != nil {
		logger.Errorf("trader: 解析json异常: %+v", err)
		return nil, err
	}
	return detail, nil
}

// CancelOrder 撤单
func CancelOrder(orderId int) error {
	params := urlpkg.Values{
		"order_id": {fmt.Sprintf("%d", orderId)},
	}
	body := params.Encode()
	logger.Infof("trader-cancel: %s", body)
	data, err := http.Post(urlCancelOrder, body)
	if err != nil {
		logger.Errorf("trader-cancel: 撤单操作异常: %+v", err)
		return err
	}
	var detail OrderResult
	err = json.Unmarshal(data, &detail)
	if err != nil {
		logger.Errorf("trader-cancel: 解析json异常: %+v", err)
		return err
	}
	logger.Infof("trader-cancel: %s, response: status=%d", body, detail.Status)
	return nil
}

// PlaceOrder 下委托订单
func PlaceOrder(direction Direction, model models.Strategy, securityCode string, price float64, volume int) (int, error) {
	_, mflag, symbol := proto.DetectMarket(securityCode)
	params := urlpkg.Values{
		"direction": {direction.String()},
		"code":      {fmt.Sprintf("%s.%s", symbol, strings.ToUpper(mflag))},
		"price":     {fmt.Sprintf("%f", price)},
		"volume":    {fmt.Sprintf("%d", volume)},
		"strategy":  {fmt.Sprintf("%d", model.Code())},
		"remark":    {model.OrderFlag()},
	}
	body := params.Encode()
	logger.Infof("trader-order: %s", body)
	data, err := http.Post(urlPlaceOrder, body)
	if err != nil {
		logger.Errorf("trader-order: 下单操作异常: %+v", err)
		return -1, err
	}
	var detail OrderResult
	err = json.Unmarshal(data, &detail)
	if err != nil {
		logger.Errorf("trader-order: 解析json异常: %+v", err)
		return -1, err
	}
	logger.Infof("trade-order: %s, response: order_id=%d", body, detail.OrderId)
	return detail.OrderId, nil
}
