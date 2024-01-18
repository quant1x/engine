package trader

import (
	"gitee.com/quant1x/exchange"
	"gitee.com/quant1x/gox/api"
	"gitee.com/quant1x/gox/logger"
	"slices"
	"sync"
)

//// HoldingOrder 持仓订单
//type HoldingOrder struct {
//	AccountType   int     `name:"账户类型" json:"account_type" dataframe:"account_type"`     // 账户类型
//	Date          string  `name:"信号日期" json:"date" dataframe:"date"`                     // 日期
//	AccountId     string  `name:"资金账户" json:"account_id" dataframe:"account_id"`         // 资金账号
//	OrderTime     string  `name:"委托时间" json:"order_time" dataframe:"order_time"`         // 报单时间
//	StockCode     string  `name:"证券代码" json:"stock_code" dataframe:"stock_code"`         // 证券代码, 例如"600000.SH"
//	OrderType     int     `name:"订单类型" json:"order_type" dataframe:"order_type"`         // 委托类型, 23:买, 24:卖
//	Price         float64 `name:"委托价格" json:"price" dataframe:"price"`                   // 报价价格, 如果price_type为指定价, 那price为指定的价格, 否则填0
//	PriceType     int     `name:"报价类型" json:"price_type" dataframe:"price_type"`         // 报价类型, 详见帮助手册
//	OrderVolume   int     `name:"委托量" json:"order_volume" dataframe:"order_volume"`      // 委托数量, 股票以'股'为单位, 债券以'张'为单位
//	OrderId       int     `name:"订单ID" json:"order_id" dataframe:"order_id"`             // 委托编号
//	OrderSysid    string  `name:"合同编号" json:"order_sysid" dataframe:"order_sysid"`       // 柜台编号
//	TradedPrice   float64 `name:"成交均价" json:"traded_price" dataframe:"traded_price"`     // 成交均价
//	TradedVolume  int     `name:"成交数量" json:"traded_volume" dataframe:"traded_volume"`   // 成交数量, 股票以'股'为单位, 债券以'张'为单位
//	OrderStatus   int     `name:"订单状态" json:"order_status" dataframe:"order_status"`     // 委托状态
//	StatusMessage string  `name:"委托状态描述" json:"status_msg" dataframe:"status_message"`   // 委托状态描述, 如废单原因
//	StrategyName  string  `name:"策略名称" json:"strategy_name" dataframe:"strategy_name"`   // 策略名称
//	OrderRemark   string  `name:"委托备注" json:"order_remark" dataframe:"order_remark"`     // 委托备注
//	HoldingPeriod int     `name:"持股周期" json:"holding_period" dataframe:"holding_period"` // 持股周期
//}
//
//// Key 索引字段: 日期/订单类型/策略名称/证券代码
//func (this HoldingOrder) Key() string {
//	return fmt.Sprintf("%s/%d/%s/%s", this.Date, this.OrderType, this.StrategyName, this.StockCode)
//}

type HoldingPosition struct {
	AccountType     int     `name:"账户类型" dataframe:"account_type"`     // 账户类型
	AccountId       string  `name:"资金账户" dataframe:"account_id"`       // 资金账号
	StockCode       string  `name:"证券代码" dataframe:"stock_code"`       // 证券代码, 例如"600000.SH"
	Volume          int     `name:"持仓数量" dataframe:"volume"`           // 持仓数量,股票以'股'为单位, 债券以'张'为单位
	CanUseVolume    int     `name:"可卖数量" dataframe:"can_use_volume"`   // 可用数量, 股票以'股'为单位, 债券以'张'为单位
	OpenPrice       float64 `name:"开仓价" dataframe:"open_price"`        // 开仓价
	MarketValue     float64 `name:"市值" dataframe:"market_value"`       // 市值
	FrozenVolume    int     `name:"冻结数量" dataframe:"frozen_volume"`    // 冻结数量
	OnRoadVolume    int     `name:"在途股份" dataframe:"on_road_volume"`   // 在途股份
	YesterdayVolume int     `name:"昨夜拥股" dataframe:"yesterday_volume"` // 昨夜拥股
	AvgPrice        float64 `name:"成本价" dataframe:"avg_price"`         // 成本价
	HoldingPeriod   int     `name:"持股周期" dataframe:"holding_period"`   // 持股周期
}

var (
	//holdingOnce   coroutine.RollingOnce
	holdingOrders []HoldingPosition
	holdingMutex  sync.RWMutex
)

// 懒加载持仓个股的持股周期
func lazyLoadHoldingOrder() {
	holdingMutex.Lock()
	defer holdingMutex.Unlock()
	methodName := "lazyLoadHoldingOrder"
	// 1. 获取持仓列表
	positions, err := QueryHolding()
	if err != nil {
		logger.Errorf("查询持仓列表失败, error=%+v", err)
		return
	}
	// 过滤掉清仓的个股
	positions = api.Filter(positions, func(detail PositionDetail) bool {
		return detail.Volume > 0
	})
	// 清空缓存
	clear(holdingOrders)
	// 2. 用持仓列表遍历历史订单缓存文件, 补全持仓订单
	for _, position := range positions {
		var holding HoldingPosition
		_ = api.Copy(&holding, &position)
		code := position.StockCode
		volume := position.Volume
		// 矫正证券代码
		securityCode := exchange.CorrectSecurityCode(position.StockCode)
		lastTradeDate := exchange.LastTradeDate()
		dates := exchange.TradingDateRange(exchange.MARKET_CH_FIRST_LISTTIME, lastTradeDate)
		// 反转日期切片
		slices.Reverse(dates)
		// 历史记录合计买数量
		tmpTradedVolume := 0
		// 最早的持股日期
		earlierDate := lastTradeDate
		// 持股周期
		holdingPeriod := 0
		// 从当前日期往前回溯订单
		for _, date := range dates {
			isLastTradeDate := date == lastTradeDate
			// 获取date的订单列表
			orders := GetOrderList(date)
			if len(orders) == 0 && isLastTradeDate {
				// 如果本地缓存订单列表为空, 且是最后一个交易日, 则从券商获取订单列表
				orders, _ = QueryOrders()
			}
			if len(orders) == 0 {
				// 如果订单列表为空, 中断循环
				break
			}
			orders = api.Filter(orders, func(detail OrderDetail) bool {
				return detail.StockCode == code
			})
			if len(orders) == 0 {
				continue
			}
			currentTradedVolume := 0
			for _, order := range orders {
				if order.OrderType != STOCK_BUY && order.OrderType != STOCK_SELL {
					// 如果不是买入和卖出, 忽略
					continue
				}
				plus := 1 // 默认相加
				if order.OrderType == STOCK_SELL {
					// 卖出则减
					plus = -1
				}
				if order.OrderStatus == ORDER_PART_SUCC || order.OrderStatus == ORDER_SUCCEEDED {
					// 部成和已成 , 成交数量累加
					currentTradedVolume += plus * order.TradedVolume
				}
			}
			earlierDate = date
			tmpTradedVolume += currentTradedVolume
			if tmpTradedVolume == volume {
				// 如果订单合计成交量等于持仓量, 则退出
				break
			}
		}
		if tmpTradedVolume != volume {
			// 历史已成买入量和持仓量不一致, 按照持仓逻辑, 会当成持股最后一天来处理
			logger.Errorf("[%s]: 加载(%s)持仓记录异常, 历史委托记录合并持仓量不一致", methodName, securityCode)
		}
		// 计算持股周期
		dates = exchange.TradingDateRange(earlierDate, lastTradeDate)
		holdingPeriod = len(dates) - 1
		holding.HoldingPeriod = holdingPeriod
		holdingOrders = append(holdingOrders, holding)
	}
	// 3. 对于非策略订单的处理
}

// GetHoldingPeriodList 获取持仓周期列表
func GetHoldingPeriodList() []HoldingPosition {
	lazyLoadHoldingOrder()
	return holdingOrders
}
