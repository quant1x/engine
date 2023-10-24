package smart

// StockPool 股票池
type StockPool struct {
	Status         int     `name:"策略状态" dataframe:"status"`
	Date           string  `name:"信号日期" dataframe:"date"`
	Code           string  `name:"证券代码" dataframe:"code"`
	Name           string  `name:"证券名称" dataframe:"name"`
	TurnZ          float64 `name:"开盘换手Z" dataframe:"turn_z"`
	Rate           float64 `name:"涨跌幅%" dataframe:"rate"`
	Buy            float64 `name:"委托价格" dataframe:"buy"`
	Sell           float64 `name:"目标价格" dataframe:"sell"`
	StrategyCode   int     `name:"策略编码" dataframe:"strategy_code"`
	StrategyName   string  `name:"策略名称" dataframe:"strategy_name"`
	Rules          uint64  `name:"规则" dataframe:"rules"`
	BlockType      string  `name:"板块类型" dataframe:"block_type"`
	BlockCode      string  `name:"板块代码" dataframe:"block_code"`
	BlockName      string  `name:"板块名称" dataframe:"block_name"`
	BlockRate      float64 `name:"板块涨幅%" dataframe:"block_rate"`
	BlockTop       int     `name:"板块排名" dataframe:"block_top"`
	BlockRank      int     `name:"个股排名" dataframe:"block_rank"`
	BlockZhangTing string  `name:"板块涨停数" dataframe:"block_zhangting"`
	BlockDescribe  string  `name:"涨/跌/平" dataframe:"block_describe"`
	BlockTopCode   string  `name:"领涨股代码" dataframe:"block_top_code"`
	BlockTopName   string  `name:"领涨股名称" dataframe:"block_top_name"`
	BlockTopRate   float64 `name:"领涨股涨幅%" dataframe:"block_top_rate"`
	Tendency       string  `name:"短线趋势" dataframe:"tendency"`
	UpdateTime     string  `name:"更新时间" dataframe:"update_time"`
}

//const (
//	RuleMiss   Kind = iota //规则未命中
//	RuleHit                // 命中
//	RuleCancel             // 撤回
//	RulePassed             // 成功
//	RuleFailed             // 失败
//)
