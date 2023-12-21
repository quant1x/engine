package config

import (
	"gitee.com/quant1x/engine/market"
	"gitee.com/quant1x/gotdx/securities"
	"gitee.com/quant1x/gox/api"
	"slices"
)

// TraderRole 交易员角色
type TraderRole int

const (
	RoleDisable TraderRole = iota // 禁止自动化交易
	RolePython                    // python脚本自动化交易
	RoleProxy                     // 代理交易模式
	RoleManual                    // 人工干预, 作用同
)

// TraderParameter 预览交易通道参数
type TraderParameter struct {
	AccountId            string         `name:"账号ID" yaml:"account_id" dataframe:"888xxxxxxx"`                                      // 账号ID
	StampDutyRateForBuy  float64        `name:"买入印花税" yaml:"stamp_duty_rate_for_buy" default:"0.0000"`                              // 印花说-买入, 没有
	StampDutyRateForSell float64        `name:"卖出印花税" yaml:"stamp_duty_rate_for_sell" default:"0.0010"`                             // 印花说-卖出, 默认是千分之1
	TransferRate         float64        `name:"过户费" yaml:"transfer_rate" default:"0.0006"`                                          // 过户费, 双向, 默认是万分之6
	CommissionRate       float64        `name:"佣金率" yaml:"commission_rate" default:"0.00025"`                                       // 券商佣金, 双向, 默认万分之2.5
	CommissionMin        float64        `name:"佣金最低" yaml:"commission_min" default:"5.0000"`                                        // 券商佣金最低, 双向, 默认5.00
	PositionRatio        float64        `name:"持仓占比" yaml:"position_ratio" default:"0.5000"`                                        // 当日持仓占比, 默认50%
	KeepCash             float64        `name:"保留现金" yaml:"keep_cash" default:"10000.00"`                                           // 保留现金, 默认10000.00
	BuyAmountMax         float64        `name:"可买最大金额" yaml:"buy_amount_max" default:"250000.00"`                                   // 买入最大金额, 默认250000.00
	BuyAmountMin         float64        `name:"可买最小金额" yaml:"buy_amount_min" default:"1000.00"`                                     // 买入最小金额, 默认1000.00
	Role                 TraderRole     `name:"角色" yaml:"role" default:"3"`                                                         // 交易员角色, 默认是需要人工干预, 系统不做自动交易处理
	ProxyUrl             string         `name:"代理URL" yaml:"proxy_url" default:"http://127.0.0.1:18168/qmt"`                        // 禁止使用公网地址
	Strategies           []TradeRule    `name:"策略集合" yaml:"strategies"`                                                             // 策略集合
	CancelSession        TradingSession `name:"撤单时段" yaml:"cancel" default:"09:15:00~09:19:59,09:25:00~11:29:59,13:00:00~14:59:59"` // 可撤单配置
	//HeadOrderAuto        bool           `name:"早盘自动买入" yaml:"head_order_auto" default:"false"`                                        // 早盘订单是否自动买入
	//TickOrderAuto        bool           `name:"盘中自动买入" yaml:"tick_order_auto" default:"false"`                                        // 盘中订单是否自动买入
	//TailOrderAuto        bool           `name:"尾盘自动买入" yaml:"tail_order_auto" default:"false"`                                        // 尾盘订单是否自动买入
	//SellOrderAuto        bool           `name:"自动卖出" yaml:"sell_order_auto" default:"false"`                                          // 是否自动卖出
}

// TotalNumberOfTargets 统计标的总数
func (t TraderParameter) TotalNumberOfTargets() int {
	total := 0
	for _, v := range t.Strategies {
		total += v.NumberOfTargets()
	}
	return total
}

// ResetPositionRatio 重置仓位占比
func (t TraderParameter) ResetPositionRatio() {
	remainingRatio := 1.00
	strategyCount := len(t.Strategies)
	var unassignedStrategies []*TradeRule
	for i := 0; i < strategyCount; i++ {
		v := &(t.Strategies[i])
		if !v.BuyEnable() {
			continue
		}
		// 校对个股最大资金
		if v.FeeMax > t.BuyAmountMax {
			v.FeeMax = t.BuyAmountMax
		}
		// 校对个股最小资金
		if v.FeeMin < t.BuyAmountMin {
			v.FeeMin = t.BuyAmountMin
		}
		if v.Weight > 1.00 {
			v.Weight = 1.00
		}
		if v.Weight > 0 {
			remainingRatio -= v.Weight
		} else {
			unassignedStrategies = append(unassignedStrategies, v)
		}
	}
	remainingCount := len(unassignedStrategies)
	if remainingRatio > 0 && remainingCount > 0 {
		averageFundPercentage := remainingRatio / float64(remainingCount)
		for _, v := range unassignedStrategies {
			v.Weight = averageFundPercentage
		}
	}
}

// TradeRule 交易规则
type TradeRule struct {
	Id                  int            `name:"策略编码" yaml:"id" default:"-1"`                                    // 策略ID, -1无效
	Auto                bool           `name:"是否自动执行" yaml:"auto" default:"false"`                             // 是否自动执行
	Name                string         `name:"策略名称" yaml:"name"`                                               // 策略名称
	Flag                string         `name:"订单标识" yaml:"flag"`                                               // 订单标识,分早盘,尾盘和盘中
	Session             TradingSession `name:"时间范围" yaml:"time" default:"09:30:00~11:30:00,13:00:00~14:56:30"` // 可操作的交易时段
	Weight              float64        `name:"持仓占比" yaml:"weight" default:"0"`                                 // 策略权重, 默认0, 由系统自动分配
	Total               int            `name:"订单数上限" yaml:"total" default:"3"`                                 // 订单总数, 默认是3
	FeeMax              float64        `name:"最大费用" yaml:"fee_max" default:"20000.00"`                         // 可投入资金-最大
	FeeMin              float64        `name:"最小费用" yaml:"fee_min" default:"10000.00"`                         // 可投入资金-最小
	Sectors             []string       `name:"板块" yaml:"sectors" default:""`                                   // 板块, 策略适用的板块列表, 默认板块为空, 即全部个股
	IgnoreMarginTrading bool           `name:"剔除两融" yaml:"ignore_margin_trading" default:"true"`               // 剔除两融标的, 默认是剔除
	HoldingPeriod       int            `name:"持仓周期" yaml:"holding_period" default:"1"`                         // 持仓周期, 默认为1天, 即T+1日触发117号策略
	SellStrategy        int            `name:"卖出策略" yaml:"sell_strategy" default:"117"`                        // 卖出策略, 默认117
	TakeProfitRatio     float64        `name:"止盈比例" yaml:"take_profit_ratio" default:"15.00"`                  // 止盈比例, 默认15%
	StopLossRatio       float64        `name:"止损比例" yaml:"stop_loss_ratio" default:"-2.00"`                    // 止损比例, 默认-2%
}

func (t *TradeRule) QmtStrategyName() string {
	return QmtStrategyNameFromId(t.Id)
}

// Enable 策略是否有效
func (t *TradeRule) Enable() bool {
	return t.Auto && t.Id >= 0
}

// BuyEnable 获取可买入状态
func (t *TradeRule) BuyEnable() bool {
	return t.Enable() && t.Total > 0
}

// SellEnable 获取可卖出状态
func (t *TradeRule) SellEnable() bool {
	return t.Enable()
}

// IsCookieCutterForSell 是否一刀切卖出
func (t *TradeRule) IsCookieCutterForSell() bool {
	return t.SellEnable() && t.Total == 0
}

// NumberOfTargets 获得可买入标的总数
func (t *TradeRule) NumberOfTargets() int {
	if !t.BuyEnable() {
		return 0
	}
	return t.Total
}

func (t *TradeRule) StockList() []string {
	var codes []string
	for _, v := range t.Sectors {
		blockInfo := securities.GetBlockInfo(v)
		if blockInfo != nil {
			codes = append(codes, blockInfo.ConstituentStocks...)
		}
	}
	if len(codes) == 0 {
		codes = market.GetCodeList()
	}
	codes = api.SliceUnique(codes, func(a string, b string) int {
		if a < b {
			return -1
		} else if a > b {
			return 1
		} else {
			return 0
		}
	})

	if t.IgnoreMarginTrading {
		// 过滤两融
		marginTradingList := securities.MarginTradingList()
		codes = api.Filter(codes, func(s string) bool {
			if slices.Contains(marginTradingList, s) {
				return false
			}
			return true
		})
	}
	return codes
}

// TraderConfig 获取交易配置
func TraderConfig() TraderParameter {
	trader := GlobalConfig.Trader
	trader.ResetPositionRatio()
	return trader
}

// GetTradeRule 通过策略编码查找规则
func GetTradeRule(code int) *TradeRule {
	strategies := TraderConfig().Strategies
	for _, v := range strategies {
		if v.Auto && v.Id == code {
			return &v
		}
	}
	return nil
}

//// GetSellRule 获取卖出规则
//func GetSellRule() TradeRule {
//	params := TraderConfig()
//	return params.Sell
//}
