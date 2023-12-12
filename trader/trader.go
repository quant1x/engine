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
	traderConfig = config.TraderConfig()
	// miniQMT代理服务器地址
	//urlPrefixMiniQmtProxy = "http://10.211.55.3:18168/qmt"
	urlPrefixMiniQmtProxy = traderConfig.ProxyUrl
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

// 交易类型标志
func (d Direction) Flag() string {
	flag := d.String()
	return flag[0:1]
}

const (
	BUY  Direction = "buy"  // 买入
	SELL Direction = "sell" // 卖出
	JUNK Direction = "junk" // 废单
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
	AccountType     int     `name:"账户类型" json:"account_type"`     // 账户类型
	AccountId       string  `name:"资金账户" json:"account_id"`       // 资金账号
	StockCode       string  `name:"证券代码" json:"stock_code"`       // 证券代码, 例如"600000.SH"
	Volume          int     `name:"持仓数量" json:"volume"`           // 持仓数量,股票以'股'为单位, 债券以'张'为单位
	CanUseVolume    int     `name:"可卖数量" json:"can_use_volume"`   // 可用数量, 股票以'股'为单位, 债券以'张'为单位
	OpenPrice       float64 `name:"开仓价" json:"open_price"`        // 开仓价
	MarketValue     float64 `name:"市值" json:"market_value"`       // 市值
	FrozenVolume    int     `name:"冻结数量" json:"frozen_volume"`    // 冻结数量
	OnRoadVolume    int     `name:"在途股份" json:"on_road_volume"`   // 在途股份
	YesterdayVolume int     `name:"昨夜拥股" json:"yesterday_volume"` // 昨夜拥股
	AvgPrice        float64 `name:"成本价" json:"avg_price"`         // 成本价
}

// OrderDetail 委托订单
type OrderDetail struct {
	AccountType   int     `name:"账户类型" json:"account_type"`  // 账户类型
	AccountId     string  `name:"资金账户" json:"account_id"`    // 资金账号
	OrderTime     string  `name:"委托时间" json:"order_time"`    // 报单时间
	StockCode     string  `name:"证券代码" json:"stock_code"`    // 证券代码, 例如"600000.SH"
	OrderType     int     `name:"订单类型" json:"order_type"`    // 委托类型, 23:买, 24:卖
	Price         float64 `name:"委托价格" json:"price"`         // 报价价格, 如果price_type为指定价, 那price为指定的价格, 否则填0
	PriceType     int     `name:"报价类型" json:"price_type"`    // 报价类型, 详见帮助手册
	OrderVolume   int     `name:"委托量" json:"order_volume"`   // 委托数量, 股票以'股'为单位, 债券以'张'为单位
	OrderId       int     `name:"订单ID" json:"order_id"`      // 委托编号
	OrderSysid    string  `name:"合同编号" json:"order_sysid"`   // 柜台编号
	TradedPrice   float64 `name:"成交均价" json:"traded_price"`  // 成交均价
	TradedVolume  int     `name:"成交数量" json:"traded_volume"` // 成交数量, 股票以'股'为单位, 债券以'张'为单位
	OrderStatus   int     `name:"订单状态" json:"order_status"`  // 委托状态
	StatusMessage string  `name:"委托状态描述" json:"status_msg"`  // 委托状态描述, 如废单原因
	StrategyName  string  `name:"策略名称" json:"strategy_name"` // 策略名称
	OrderRemark   string  `name:"委托备注" json:"order_remark"`  // 委托备注
}

// SecurityCode 获取证券代码
func (d OrderDetail) SecurityCode() string {
	return proto.CorrectSecurityCode(d.StockCode)
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
	strategyName := models.QmtStrategyName(model)
	orderRemark := models.QmtOrderRemark(model)
	return DirectOrder(direction, strategyName, orderRemark, securityCode, price, volume)
}

// 直接下单(透传)
func DirectOrder(direction Direction, strategyName, orderRemark, securityCode string, price float64, volume int) (int, error) {
	_, mflag, symbol := proto.DetectMarket(securityCode)
	params := urlpkg.Values{
		"direction": {direction.String()},
		"code":      {fmt.Sprintf("%s.%s", symbol, strings.ToUpper(mflag))},
		"price":     {fmt.Sprintf("%f", price)},
		"volume":    {fmt.Sprintf("%d", volume)},
		"strategy":  {strategyName},
		"remark":    {orderRemark},
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

// 计算策略标的的可用资金
func CalculateFundForStrategy(model models.Strategy) float64 {
	strategyCode := model.Code()
	rule := config.GetTradeRule(strategyCode)
	if rule == nil {
		return InvalidFee
	}
	fund := CalculateAvailableFund(rule)
	return fund
}
