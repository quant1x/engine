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
	Role             TraderRole     `name:"角色" yaml:"role" default:"3"`                                  // 交易员角色
	ProxyUrl         string         `name:"代理URL" yaml:"proxy_url" default:"http://127.0.0.1:18168/qmt"` // 禁止使用公网地址
	Head             TradeRule      `name:"早盘" yaml:"head" default:"{\"Time\":\"09:30:00~11:30:00\"}"`
	Tail             TradeRule      `name:"尾盘" yaml:"tail" default:"{\"Time\":\"14:50:00~14:56:30\"}"`
	Tick             TradeRule      `name:"盘中" yaml:"tick" default:"{\"Time\":\"09:39:00~14:56:30\"}"`
	Sell             TradeRule      `name:"卖出" yaml:"sell" default:""`
	ReservedOfCancel string         `name:"撤单保留字段" yaml:"cancel" default:"09:15:00~09:19:59,09:25:00~11:29:59,13:00:00~14:59:59"`
	CancelSession    TradingSession `name:"撤单时段" yaml:"-" default:""`
}

// TradeRule 交易规则
type TradeRule struct {
	Auto    bool           `name:"是否自动执行" yaml:"auto" default:"false"`
	Time    string         `name:"时间范围" yaml:"time" default:"09:30:00~11:30:00,13:00:00~14:56:30"`
	Session TradingSession `name:"交易时段" yaml:"-"`
	Total   int            `name:"订单数上限" yaml:"total" default:"3"`
	FeeMax  float64        `name:"最大费用" yaml:"fee_max" default:"20000.00"`
	FeeMin  float64        `name:"最小费用" yaml:"fee_min" default:"10000.00"`
}

// TraderConfig 获取交易配置
func TraderConfig() TraderParameter {
	return GlobalConfig.Trader
}
