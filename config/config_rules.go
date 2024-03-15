package config

// RuleParameter 规则参数
type RuleParameter struct {
	SectorsFilter               bool        `yaml:"sectors_filter" default:"false"`              // 是否启用板块过滤, false代表全市场扫描
	SectorsTopN                 int         `yaml:"sectors_top_n" default:"3"`                   // 最多关联多少个板块, 默认3个
	StockTopNInSector           int         `yaml:"stock_top_n_in_sector" default:"5"`           // 板块内个股排名前N
	IgnoreRuleGroup             []int       `yaml:"ignore_rule_group"`                           // 忽略规则组合
	IgnoreCodes                 []string    `yaml:"ignore_codes" default:"[\"sh68\",\"bj\"]"`    // 忽略的证券代码段, 默认忽略科创板和北交所全部
	MaximumIncreaseWithin5days  float64     `yaml:"maximum_increase_within_5d" default:"20.00"`  // 20.00 5日累计最大涨幅
	MaximumIncreaseWithin10days float64     `yaml:"maximum_increase_within_10d" default:"70.00"` // 70.00 10日累计最大涨幅
	MaxReduceAmount             float64     `yaml:"max_reduce_amount" default:"-1000"`           // -1000 最大流出1000万
	SafetyScore                 NumberRange `yaml:"safety_score" default:"80~"`                  // 80 通达信安全分最小值
	VolumeRatio                 NumberRange `yaml:"volume_ratio" default:"0.382~2.800"`          // 1.800 成交量放大不能超过1.8
	Capital                     NumberRange `yaml:"capital" default:"0.5~20"`                    // 流通股本, 默认0.5亿~20亿
	MarketCap                   NumberRange `yaml:"market_cap" default:"4~600"`                  // 流通市值, 默认4亿~600亿
	Price                       NumberRange `yaml:"price" default:"2~"`                          // 股价: 4.9E-324~1.7976931348623157e+308
	OpenChangeRate              NumberRange `yaml:"open_change_rate"  default:""`                // 开盘涨幅, 默认不限制
	OpenQuantityRatio           NumberRange `yaml:"open_quantity_ratio" default:""`              // 开盘量比, 默认不限制
	OpenTurnZ                   NumberRange `yaml:"open_turn_z" default:""`                      // 开盘换手, 默认不限制
	ChangeRate                  NumberRange `yaml:"change_rate"  default:""`                     // 涨幅, 默认不限制
	Vix                         NumberRange `yaml:"vix" default:""`                              // 波动率, 默认不限制
	TurnoverRate                NumberRange `yaml:"turnover_rate" default:""`                    // 换手率范围, 默认不限制
	AmplitudeRatio              NumberRange `yaml:"amplitude_ratio" default:""`                  // 振幅范围, 默认不限制
	BiddingVolume               NumberRange `yaml:"bidding_volume" default:""`                   // 5档行情委托平均值范围, 默认不限制
	Sentiment                   NumberRange `yaml:"sentiment" default:"38.2~61.80"`              // 情绪范围
	GapDown                     bool        `yaml:"gap_down" default:"true"`                     // 买入是否允许跳空低开, 默认是允许
}
