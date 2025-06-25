package trader

import (
	"gitee.com/quant1x/exchange"
	"gitee.com/quant1x/gox/api"
	"gitee.com/quant1x/gox/logger"
	"slices"
	"sync"
)

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
	dates := GetLocalOrderDates()
	if len(dates) == 0 {
		return
	}
	// 3. 重新评估持仓范围, 有可能存在日期没有成交的可能
	firstDate := dates[0]
	lastTradeDate := exchange.LastTradeDate()
	dates = exchange.TradingDateRange(firstDate, lastTradeDate)
	// 反转日期切片
	slices.Reverse(dates)
	// 4. 用持仓列表遍历历史订单缓存文件, 补全持仓订单
	for _, position := range positions {
		var holding HoldingPosition
		_ = api.Copy(&holding, &position)
		code := position.StockCode
		volume := position.Volume
		// 矫正证券代码
		securityCode := exchange.CorrectSecurityCode(position.StockCode)
		// 历史记录合计买数量
		tmpTradedVolume := 0
		// 最早的持股日期
		earlierDate := lastTradeDate
		// 持股周期
		holdingPeriod := 0
		//if code == "001216.SZ" {
		//	fmt.Println(code)
		//}
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
				// 如果订单列表为空, 跳过
				continue
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
		dateRanges := exchange.TradingDateRange(earlierDate, lastTradeDate)
		holdingPeriod = len(dateRanges) - 1
		holding.HoldingPeriod = holdingPeriod
		holdingOrders = append(holdingOrders, holding)
	}
	// 5. 对于非策略订单的处理
}

// GetHoldingPeriodList 获取持仓周期列表
func GetHoldingPeriodList() []HoldingPosition {
	lazyLoadHoldingOrder()
	return holdingOrders
}
