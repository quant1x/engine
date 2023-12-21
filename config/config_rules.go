package config

// RuleParameter 规则参数
type RuleParameter struct {
	SectorsFilter               bool        `yaml:"sectors_filter" default:"false"`                  // 是否启用板块过滤, false代表全市场扫描
	SectorsTopN                 int         `yaml:"sectors_top_n" default:"3"`                       // 最多关联多少个板块, 默认3个
	StockTopNInSector           int         `yaml:"stock_top_n_in_sector" default:"5"`               // 板块内个股排名前N
	IgnoreCodes                 []string    `yaml:"ignore_codes" default:"[\"sh68\",\"bj\"]"`        // 忽略的证券代码段, 默认忽略科创板和北交所全部
	SafetyScoreMin              float64     `yaml:"safety_score_min" default:"80"`                   // 80 通达信安全分最小值
	MaximumIncreaseWithin5days  float64     `yaml:"maximum_increase_within_5_days" default:"20.00"`  // 20.00 5日累计最大涨幅
	MaximumIncreaseWithin10days float64     `yaml:"maximum_increase_within_10_days" default:"70.00"` // 70.00 10日累计最大涨幅
	MaxReduceAmount             float64     `yaml:"max_reduce_amount" default:"-1000"`               // -1000 最大流出1000万
	VolumeRatioMax              float64     `yaml:"volume_ratio_max" default:"3.82"`                 // 1.800 成交量放大不能超过1.8
	Price                       NumberRange `yaml:"price" default:"2.00~30.00"`                      // 股价
	Capital                     NumberRange `yaml:"capital" default:"2~20"`                          // 流通股本
	MarketCap                   NumberRange `yaml:"market_cap" default:"4~600"`                      // 流通市值
	OpenRate                    NumberRange `yaml:"open_rate"  default:"-2.00~2.00"`                 // 开盘涨幅
	QuantityRatio               NumberRange `yaml:"quantity_ratio" default:"1.00~9.99"`              // 开盘量比
	TurnZ                       NumberRange `yaml:"turn_z" default:"1.50~200.00"`                    // 开盘换手
	Vix                         NumberRange `yaml:"vix" default:"0.00~100.00"`                       // 波动率0
	TurnoverRate                NumberRange `yaml:"turnover_rate" default:"1.00~20.00"`              // 换手率范围
	AmplitudeRatio              NumberRange `yaml:"amplitude_ratio" default:"0.00~15.00"`            // 振幅 最大
	BiddingVolume               NumberRange `yaml:"bidding_volume" default:"100~5000"`               // 5档行情委托平均最小值
	Sentiment                   NumberRange `yaml:"sentiment" default:"38.2~61.80"`                  // 情绪范围
}

// RuleConfig 获取交易配置
func RuleConfig() RuleParameter {
	return GlobalConfig.Rules
}
