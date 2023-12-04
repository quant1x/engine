package config

// PreviewTraderParameter 预览交易通道参数
type PreviewTraderParameter struct {
	ProxyUrl string    `name:"代理URL" yaml:"proxy_url" default:"http://127.0.0.1:18168/qmt"` // 禁止使用公网地址
	Head     TradeRule `name:"早盘" yaml:"head" default:"{\"Time\":\"09:30:00~11:30:00\"}"`
	Tail     TradeRule `name:"尾盘" yaml:"tail" default:"{\"Time\":\"14:50:00~14:56:30\"}"`
	Tick     TradeRule `name:"盘中" yaml:"tick" default:""`
	Sell     TradeRule `name:"卖出" yaml:"sell" default:""`
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
