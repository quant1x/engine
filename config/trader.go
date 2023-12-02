package config

// PreviewTraderParameter 预览交易通道参数
type v1PreviewTraderParameter struct {
	HeadOrderAuto      bool    `yaml:"head_order_auto" default:"true"`           // head订单自动买入
	TickOrderAuto      bool    `yaml:"tick_order_auto" default:"true"`           // tick订单自动买入
	TailOrderAuto      bool    `yaml:"tail_order_auto" default:"true"`           // tail订单自动买入
	SellOrderAuto      bool    `yaml:"sell_order_auto" default:"true"`           // 是否自动卖出
	HeadTime           string  `yaml:"head_time" default:"09:27~14:59:30"`       // head订单时间
	TailTime           string  `yaml:"tail_time" default:"14:45~14:59:30"`       // head订单时间
	TickTime           string  `yaml:"tick_time" default:"09:39:00~14:56:30"`    // 盘中实时买入tail订单时段
	AskTime            string  `yaml:"ask_time" default:"09:50~14:59:30"`        // 卖出时段
	TickOrderMaxAmount float64 `yaml:"tick_order_max_amount" default:"20000.00"` // tick订单最大买入金额
	TickOrderMinAmount float64 `yaml:"tick_order_min_amount" default:"10000.00"` // tick订单最小买入金额
	BuyAmountMax       float64 `yaml:"buy_amount_max" default:"20000.00"`        // 早盘最大可买金额
	BuyAmountMin       float64 `yaml:"buy_amount_min" default:"10000.00"`        // 早盘最小可买金额
}

// PreviewTraderParameter 预览交易通道参数
type PreviewTraderParameter struct {
	Head TradeRule `name:"早盘" yaml:"head" default:"{\"Time\":\"09:30:00~11:30:00\"}"`
	Tail TradeRule `name:"尾盘" yaml:"tail" default:"{\"Time\":\"14:50:00~14:56:30\"}"`
	Tick TradeRule `name:"盘中" yaml:"tick" default:""`
	Sell TradeRule `name:"卖出" yaml:"sell" default:""`
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
