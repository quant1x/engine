package config

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
	PositionRatio        float64        `name:"持仓占比" yaml:"position_ratio" default:"0.5000"`                                          // 当日持仓占比, 默认50%
	StampDutyRateForBuy  float64        `name:"买入印花税" yaml:"stamp_duty_rate_for_buy" default:"0.0000"`                                // 印花说-买入, 没有
	StampDutyRateForSell float64        `name:"卖出印花税" yaml:"stamp_duty_rate_for_sell" default:"0.0010"`                               // 印花说-卖出, 默认是千分之1
	TransferRate         float64        `name:"过户费" yaml:"transfer_rate" default:"0.0006"`                                            // 过户费, 双向, 默认是万分之6
	CommissionRate       float64        `name:"佣金率" yaml:"commission_rate" default:"0.00025"`                                         // 券商佣金, 双向, 默认万分之2.5
	CommissionMin        float64        `name:"佣金最低" yaml:"commission_min" default:"5.0000"`                                          // 券商佣金最低, 双向, 默认5.00
	KeepCash             float64        `name:"保留现金" yaml:"keep_cash" default:"10000.00"`                                             // 保留现金, 默认10000.00
	Role                 TraderRole     `name:"角色" yaml:"role" default:"3"`                                                           // 交易员角色, 默认是需要人工干预, 系统不做自动交易处理
	ProxyUrl             string         `name:"代理URL" yaml:"proxy_url" default:"http://127.0.0.1:18168/qmt"`                          // 禁止使用公网地址
	Head                 TradeRule      `name:"早盘" yaml:"head" default:"{\"Time\":\"09:30:00~11:30:00\"}"`                            // 订单策略配置-早盘
	Tail                 TradeRule      `name:"尾盘" yaml:"tail" default:"{\"Time\":\"14:50:00~14:56:30\"}"`                            // 订单策略配置-尾盘
	Tick                 TradeRule      `name:"盘中" yaml:"tick" default:"{\"Time\":\"09:39:00~14:56:30\"}"`                            // 订单策略配置-盘中
	Sell                 TradeRule      `name:"卖出" yaml:"sell" default:""`                                                            // 卖出策略配置
	ReservedOfCancel     string         `name:"撤单保留字段" yaml:"cancel" default:"09:15:00~09:19:59,09:25:00~11:29:59,13:00:00~14:59:59"` // 预览-可撤单配置
	CancelSession        TradingSession `name:"撤单时段" yaml:"-" default:""`                                                             // 可撤单配置
}

// TradeRule 交易规则
type TradeRule struct {
	Auto    bool           `name:"是否自动执行" yaml:"auto" default:"false"`                             // 是否自动执行
	Time    string         `name:"时间范围" yaml:"time" default:"09:30:00~11:30:00,13:00:00~14:56:30"` // 预览-执行操作的时间段
	Session TradingSession `name:"交易时段" yaml:"-"`                                                  // 可操作的交易时段
	Total   int            `name:"订单数上限" yaml:"total" default:"3"`                                 // 订单总数, 默认是3
	FeeMax  float64        `name:"最大费用" yaml:"fee_max" default:"20000.00"`                         // 可投入资金-最大
	FeeMin  float64        `name:"最小费用" yaml:"fee_min" default:"10000.00"`                         // 可投入资金-最小
}

// TraderConfig 获取交易配置
func TraderConfig() TraderParameter {
	return GlobalConfig.Trader
}
