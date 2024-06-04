package config

import "gitee.com/quant1x/exchange"

// TraderRole 交易员角色
type TraderRole int

const (
	RoleDisable TraderRole = iota // 禁止自动化交易
	RolePython                    // python脚本自动化交易
	RoleProxy                     // 代理交易模式
	RoleManual                    // 人工干预, 作用同
)

const (
	sectorIgnorePrefix = "-"
	sectorPrefixLength = len(sectorIgnorePrefix)
)

// TraderParameter 预览交易通道参数
type TraderParameter struct {
	AccountId                   string              `name:"账号ID" yaml:"account_id" dataframe:"888xxxxxxx"`                                      // 账号ID
	OrderPath                   string              `name:"订单路径" yaml:"order_path"`                                                             // 订单路径
	TopN                        int                 `name:"TopN" yaml:"top_n" default:"3"`                                                      // 最多输出前多少名个股
	HaveETF                     bool                `name:"是否包含ETF" yaml:"have_etf" default:"false"`                                            // 是否包含ETF
	PriceCageRatio              float64             `name:"价格笼子比例" yaml:"price_cage_ratio" default:"0.02"`                                      // 价格笼子比例, 默认2%, 小于0就是无限制
	MinimumPriceFluctuationUnit float64             `name:"价格变动最小单位" yaml:"minimum_price_fluctuation_unit" default:"0.10"`                      // 价格最小变动单位, 默认0.10
	AnnualInterestRate          float64             `name:"年利率" yaml:"annual_interest_rate" default:"1.65"`                                     // 2024年2月18日建设银行1年期存款利率1.65%
	StampDutyRateForBuy         float64             `name:"买入印花税" yaml:"stamp_duty_rate_for_buy" default:"0.0000"`                              // 印花说-买入, 没有
	StampDutyRateForSell        float64             `name:"卖出印花税" yaml:"stamp_duty_rate_for_sell" default:"0.0010"`                             // 印花说-卖出, 默认是千分之1
	TransferRate                float64             `name:"过户费" yaml:"transfer_rate" default:"0.0006"`                                          // 过户费, 双向, 默认是万分之6
	CommissionRate              float64             `name:"佣金率" yaml:"commission_rate" default:"0.00025"`                                       // 券商佣金, 双向, 默认万分之2.5
	CommissionMin               float64             `name:"佣金最低" yaml:"commission_min" default:"5.0000"`                                        // 券商佣金最低, 双向, 默认5.00
	PositionRatio               float64             `name:"持仓占比" yaml:"position_ratio" default:"0.5000"`                                        // 当日持仓占比, 默认50%
	KeepCash                    float64             `name:"保留现金" yaml:"keep_cash" default:"10000.00"`                                           // 保留现金, 默认10000.00
	BuyAmountMax                float64             `name:"可买最大金额" yaml:"buy_amount_max" default:"250000.00"`                                   // 买入最大金额, 默认250000.00
	BuyAmountMin                float64             `name:"可买最小金额" yaml:"buy_amount_min" default:"1000.00"`                                     // 买入最小金额, 默认1000.00
	Role                        TraderRole          `name:"角色" yaml:"role" default:"3"`                                                         // 交易员角色, 默认是需要人工干预, 系统不做自动交易处理
	ProxyUrl                    string              `name:"代理URL" yaml:"proxy_url" default:"http://127.0.0.1:18168/qmt"`                        // 禁止使用公网地址
	Strategies                  []StrategyParameter `name:"策略集合" yaml:"strategies"`                                                             // 策略集合
	CancelSession               TradingSession      `name:"撤单时段" yaml:"cancel" default:"09:15:00~09:19:59,09:25:00~11:29:59,13:00:00~14:59:59"` // 可撤单配置
	UndertakeRatio              float64             `name:"承接比" yaml:"undertake_ratio" default:"0.8000"`
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
	var unassignedStrategies []*StrategyParameter
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
		for i, v := range unassignedStrategies {
			v.Weight = averageFundPercentage
			if i+1 == remainingCount {
				v.Weight = remainingRatio
			}
			remainingRatio -= v.Weight
		}
	}
}

// DailyRiskFreeRate 计算每日无风险利率
func (t TraderParameter) DailyRiskFreeRate(date string) float64 {
	date = exchange.FixTradeDate(date)
	year := date[0:4]
	start := year + "-01-01"
	end := year + "-12-31"
	dates := exchange.DateRange(start, end)
	count := len(dates)
	return t.AnnualInterestRate / float64(count)
}

// TraderConfig 获取交易配置
func TraderConfig() TraderParameter {
	trader := GlobalConfig.Trader
	trader.ResetPositionRatio()
	return trader
}

func GetSectorIgnorePrefix() string {
	return sectorIgnorePrefix
}

// GetStrategyParameterByCode 通过策略编码查找规则
func GetStrategyParameterByCode(strategyCode uint64) *StrategyParameter {
	strategies := TraderConfig().Strategies
	for _, v := range strategies {
		if v.Auto && v.Id == strategyCode {
			return &v
		}
	}
	return nil
}
